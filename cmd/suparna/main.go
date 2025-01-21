package main

import (
	"log"
	"suparna/internal/database"
	"suparna/internal/ui"
)

func main() {
	log.Println("Starting Suparna application...")
	// Initialize database
	err := database.Initialize("suparna.db")
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	log.Println("Database initialized")
	defer database.Close()

	// Example usage of filesystem for setup
	// rootPath := "/path/to/scan"
	// _, fsErr := filesystem.ScanDirectory(rootPath)
	// if fsErr != nil {
	// 	log.Printf("Error during initial scan: %v", fsErr)
	// }

	// Launch the UI
	ui.Launch()
}
