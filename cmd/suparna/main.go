package main

import (
	"log"
	"suparna/internal/database"
	"suparna/internal/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	log.Println("Starting Suparna application...")
	// Initialize the app with a unique ID
	myApp := app.NewWithID("suparna")
	window := myApp.NewWindow("Suparna File Manager")

	// Initialize the database
	err := database.Initialize("suparna.db")
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Close()

	// Launch the UI, passing the window instance
	ui.Launch(window)
}
