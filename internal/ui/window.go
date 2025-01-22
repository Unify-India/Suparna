package ui

import (
	"image/color"
	"log"
	"strconv"
	"suparna/internal/filesystem"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Set styled text with proper error or success styling
func setStyledText(output *canvas.Text, message string, isError bool) {
	output.Text = message
	if isError {
		output.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red for errors
	} else {
		bgColor := fyne.CurrentApp().Settings().Theme().Color(theme.ColorNameBackground, theme.VariantDark)
		if isDarkTheme(bgColor) {
			output.Color = color.White
		} else {
			output.Color = color.Black
		}
	}
	output.Refresh()
}

// Helper function to check if the theme is dark
func isDarkTheme(bgColor color.Color) bool {
	r, g, b, _ := bgColor.RGBA()
	return (r+g+b)/3 < 128*256
}

func Launch() {
	log.Println("Starting Suparna application...")

	myApp := app.NewWithID("suparna")
	window := myApp.NewWindow("Suparna File Manager")
	window.Resize(fyne.NewSize(1000, 800))

	// UI Components
	label := widget.NewLabel("Welcome to Suparna!")
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter directory path...")

	output := canvas.NewText("", color.Black)

	// Directory selection button
	selectButton := widget.NewButton("Select Directory", func() {
		dialog.NewFolderOpen(func(uc fyne.ListableURI, err error) {
			if err == nil && uc != nil {
				entry.SetText(uc.Path())
			}
		}, window).Show()
	})

	// File data and table initialization
	fileData := [][]string{}
	fileTable := widget.NewTable(
		func() (int, int) {
			return len(fileData), 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row < len(fileData) {
				label.SetText(fileData[id.Row][id.Col])
			}
		},
	)

	// Scan button to process the directory
	scanButton := widget.NewButton("Scan Directory", func() {
		dirPath := entry.Text
		if dirPath == "" {
			setStyledText(output, "Please select a directory first.", true)
			return
		}

		files, err := filesystem.ScanDirectoryAndSaveMetadata(dirPath)
		if err != nil {
			setStyledText(output, "Error: "+err.Error(), true)
			return
		}

		setStyledText(output, "Directory scanned successfully.", false)

		fileData = [][]string{} // Clear existing data
		for _, file := range files {
			fileData = append(fileData, []string{
				file.Name,
				file.Path,
				strconv.FormatInt(file.Size, 10),
				file.ModifiedTime.Format(time.RFC1123),
				file.Hash,
			})
		}
		fileTable.Refresh()
	})

	// Table scroll container
	tableContainer := container.NewScroll(fileTable)
	tableContainer.SetMinSize(fyne.NewSize(900, 300))

	// Layout and content setup
	content := container.NewVBox(
		label,
		entry,
		selectButton,
		scanButton,
		output,
		tableContainer,
	)
	window.SetContent(content)
	window.ShowAndRun()
}
