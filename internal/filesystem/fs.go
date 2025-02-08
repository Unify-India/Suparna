package filesystem

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"suparna/internal/database"
	"sync"
	"time"
)

// FileMetadata struct
type FileMetadata struct {
	Name         string
	Path         string
	Size         int64
	ModifiedTime time.Time
	Hash         string
}

var scanningAborted bool
var scanningMutex sync.Mutex

// getDirectoryCount returns the number of files in the DB whose path starts with the given root.
func getDirectoryCount(root string) (int, error) {
	db := database.GetDB()
	// Ensure that root ends with a separator for more accurate matching.
	pattern := root
	if pattern[len(pattern)-1] != os.PathSeparator {
		pattern += string(os.PathSeparator)
	}
	pattern += "%" // Match all files under the given directory
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM files WHERE path LIKE ?", pattern).Scan(&count)
	return count, err
}

// ScanDirectory scans a directory and reports progress via callback
func ScanDirectory(root string, progressCallback func(currentFile string, progress float64)) error {
	// Check if the directory is already crawled; if yes, skip re scanning.
	if count, err := getDirectoryCount(root); err == nil && count > 0 {
		// Optionally inform the UI that scanning is complete.
		progressCallback("Directory already scanned", 1.0)
		return nil
	}

	fileChan := make(chan FileMetadata, 100) // Buffered channel for concurrent processing

	// Reset abort flag before scan
	scanningMutex.Lock()
	scanningAborted = false
	scanningMutex.Unlock()

	// Get total file count first (for progress calculation)
	totalFiles := countFiles(root)
	if totalFiles == 0 {
		return nil
	}

	// Start file processing in a goroutine
	go func() {
		fileIndex := 0
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			// Check if scan was aborted
			scanningMutex.Lock()
			if scanningAborted {
				scanningMutex.Unlock()
				return filepath.SkipDir
			}
			scanningMutex.Unlock()

			if err != nil {
				return nil // Skip errors
			}
			if !d.IsDir() {
				info, _ := d.Info()
				hash, _ := computeHash(path)
				fileChan <- FileMetadata{
					Name:         info.Name(),
					Path:         path,
					Size:         info.Size(),
					ModifiedTime: info.ModTime(),
					Hash:         hash,
				}

				// Report progress
				fileIndex++
				progressCallback(info.Name(), float64(fileIndex)/float64(totalFiles))
			}
			return nil
		})
		close(fileChan)
	}()

	// Insert into database
	db := database.GetDB()
	tx, _ := db.Begin() // Start transaction

	for file := range fileChan {
		// Check for duplicate by file name or hash
		var originalID int64
		err := tx.QueryRow("SELECT id FROM files WHERE name=? OR hash=?", file.Name, file.Hash).Scan(&originalID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Not found, insert into files table
				res, err := tx.Exec(`INSERT INTO files (name, path, size, modified_time, hash)
					VALUES (?, ?, ?, ?, ?)`, file.Name, file.Path, file.Size, file.ModifiedTime, file.Hash)
				if err != nil {
					// ...existing error handling...
					log.Println("Error inserting into files:", err)
					continue
				}
				originalID, _ = res.LastInsertId()
			} else {
				log.Println("Error during duplicate check:", err)
				continue
			}
		} else {
			// Duplicate found, insert into duplicates table
			_, err = tx.Exec(`INSERT INTO duplicates (original_file_id, name, path, size, modified_time, hash)
				VALUES (?, ?, ?, ?, ?, ?)`, originalID, file.Name, file.Path, file.Size, file.ModifiedTime, file.Hash)
			if err != nil {
				log.Println("Error inserting into duplicates:", err)
			}
		}
	}

	tx.Commit() // Commit all inserts in one go
	return nil
}

// countFiles gets the total number of files in the directory (for progress tracking)
func countFiles(root string) int {
	count := 0
	_ = filepath.WalkDir(root, func(_ string, d os.DirEntry, _ error) error {
		if !d.IsDir() {
			count++
		}
		return nil
	})
	return count
}

// AbortScan allows stopping the scan process
func AbortScan() {
	scanningMutex.Lock()
	scanningAborted = true
	scanningMutex.Unlock()
}

// computeHash computes the MD5 hash of a file without loading it fully into RAM
func computeHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	buffer := make([]byte, 4096) // 4 KB buffer

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			hash.Write(buffer[:n]) // Process chunk
		}
		if err != nil {
			break // Stop reading on EOF
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
