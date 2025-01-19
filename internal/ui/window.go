package ui

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewApp creates a new Fyne application window.
func NewApp() *app.App {
	a := app.New()
	w := a.NewWindow("Suparna File Manager")

	// Create a simple label
	hello := widget.NewLabel("Hello, Suparna!")

	// Set the window content
	w.SetContent(container.NewVBox(
		hello,
	))

	w.ShowAndRun()
	return a
}
