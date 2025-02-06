package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Launch initializes the main application window with tabs
func Launch(window fyne.Window) {
	// Define tabs
	scanTab := CreateScanTab()
	browseTab := widget.NewLabel("Browse functionality coming soon...")
	healthTab := widget.NewLabel("Health Check functionality coming soon...")

	// Create a tab container
	tabs := container.NewAppTabs(
		container.NewTabItem("Scan", scanTab),
		container.NewTabItem("Browse", browseTab),
		container.NewTabItem("Health", healthTab),
	)

	// Set tabs as content and show the window
	window.SetContent(tabs)
	window.Resize(fyne.NewSize(1000, 800))
	window.ShowAndRun()
}
