package filesystem

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"suparna/internal/database"
	"time"
)

// FileMetadata represents metadata for a file
type FileMetadata struct {
	Name         string
	Path         string
	Size         int64
	ModifiedTime time.Time
	Hash         string
}

// ScanDirectoryAndSaveMetadata scans a directory, saves metadata to SQLite, and returns the file metadata
func ScanDirectoryAndSaveMetadata(root string) ([]FileMetadata, error) {
	var files []FileMetadata
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error reading file %s: %v\n", path, err)
			return nil // Skip this file and continue with the next
		}
		if !info.IsDir() {
			hash, hashErr := computeHash(path)
			if hashErr != nil {
				log.Printf("Error computing hash for %s: %v\n", path, hashErr)
				return nil // Skip this file and continue with the next
			}

			file := FileMetadata{
				Name:         info.Name(),
				Path:         path,
				Size:         info.Size(),
				ModifiedTime: info.ModTime(),
				Hash:         hash,
			}
			files = append(files, file)

			// Insert metadata into the database
			db := database.GetDB()
			_, insertErr := db.Exec(`INSERT OR REPLACE INTO files (name, path, size, modified_time, hash) 
				VALUES (?, ?, ?, ?, ?)`,
				file.Name, file.Path, file.Size, file.ModifiedTime, file.Hash)
			if insertErr != nil {
				log.Printf("Error inserting file %s into database: %v\n", path, insertErr)
			}
		}
		return nil
	})

	// Log all files added to the database
	// log.Printf("Files scanned and inserted: %+v", files)

	return files, err
}

// computeHash computes the MD5 hash of a file
func computeHash(path string) (string, error) {
	// log.Printf("Computing hash for %s\n", path)
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	_, err = os.Stat(path)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 4096)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			hash.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
