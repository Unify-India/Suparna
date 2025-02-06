package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CreateBrowseTab defines the "Browse Indexed Files" tab.
func CreateBrowseTab() fyne.CanvasObject {
	// Placeholder content
	placeholder := widget.NewLabel("Browse Indexed Files functionality goes here.")
	return container.NewVBox(placeholder)
}
