package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// WalkDir walks through the given directory and its subdirectories.
func WalkDir(dir string) error {
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}

		if d.IsDir() {
			// Handle directories here (e.g., print directory name)
			fmt.Println("Directory:", path)
		} else {
			// Handle files here (e.g., print file name and collect metadata)
			fmt.Println("File:", path)
			// Add code to collect file metadata (size, creation time, etc.)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}
	return nil
}

// GetFileInfo gets basic information about a file.
func GetFileInfo(filePath string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %w", err)
	}
	return fileInfo, nil
}

// (Add more functions for file system operations as needed)
