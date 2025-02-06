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

// FileMetadata struct
type FileMetadata struct {
	Name         string
	Path         string
	Size         int64
	ModifiedTime time.Time
	Hash         string
}

// ScanDirectoryAndSaveMetadata scans directory and saves metadata
func ScanDirectoryAndSaveMetadata(root string) ([]FileMetadata, error) {
	var files []FileMetadata
	fileChan := make(chan FileMetadata, 100) // Buffered channel for concurrent processing

	go func() {
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
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
			}
			return nil
		})
		close(fileChan)
	}()

	db := database.GetDB()
	tx, _ := db.Begin() // Start transaction

	for file := range fileChan {
		files = append(files, file)
		_, err := tx.Exec(`INSERT OR REPLACE INTO files (name, path, size, modified_time, hash) 
			VALUES (?, ?, ?, ?, ?)`, file.Name, file.Path, file.Size, file.ModifiedTime, file.Hash)
		if err != nil {
			log.Println("Error inserting:", err)
		}
	}

	tx.Commit() // Commit all inserts in one go
	return files, nil
}

// Compute hash
func computeHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	_, _ = file.WriteTo(hash)
	return hex.EncodeToString(hash.Sum(nil)), nil
}
