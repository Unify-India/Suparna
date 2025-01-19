package filesystem

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
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

// ScanDirectory scans a directory and returns file metadata
func ScanDirectory(root string) ([]FileMetadata, error) {
	var files []FileMetadata
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			hash, _ := computeHash(path)
			files = append(files, FileMetadata{
				Name:         info.Name(),
				Path:         path,
				Size:         info.Size(),
				ModifiedTime: info.ModTime(),
				Hash:         hash,
			})
		}
		return nil
	})
	return files, err
}

// computeHash computes the MD5 hash of a file
func computeHash(path string) (string, error) {
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
