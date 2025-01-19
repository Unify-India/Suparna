package ui

import (
	"suparna/internal/filesystem"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Launch() {
	myApp := app.New()
	window := myApp.NewWindow("Suparna File Manager")

	label := widget.NewLabel("Welcome to Suparna!")
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter directory path...")

	output := widget.NewMultiLineEntry()
	output.Disable() // Makes the MultiLineEntry read-only

	button := widget.NewButton("Scan Directory", func() {
		files, err := filesystem.ScanDirectory(entry.Text)
		if err != nil {
			output.SetText("Error: " + err.Error())
			return
		}
		output.SetText("") // Clear the output before appending new text
		for _, file := range files {
			output.SetText(output.Text + file.Path + "\n")
		}
	})

	content := container.NewVBox(label, entry, button, output)
	window.SetContent(content)
	window.ShowAndRun()
}
