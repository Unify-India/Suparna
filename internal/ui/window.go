package ui

import (
	"image/color"
	"log"
	"strconv"
	"strings"
	"suparna/internal/filesystem"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// FileData represents a single row in the UI table
type FileData struct {
	Name         string
	Path         string
	Size         string
	ModifiedTime string
	Hash         string
}

var files []FileData // Global slice to hold file data

// UI Function
func Launch() {
	log.Println("Starting Suparna application...")

	myApp := app.NewWithID("suparna")
	window := myApp.NewWindow("Suparna File Manager")
	window.Resize(fyne.NewSize(1000, 800))

	label := widget.NewLabel("Welcome to Suparna!")
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter directory path...")

	output := canvas.NewText("", color.Black)
	spinner := widget.NewProgressBarInfinite() // Progress indicator
	spinner.Hide()

	// Directory selection button
	selectButton := widget.NewButton("Select Directory", func() {
		dialog.NewFolderOpen(func(uc fyne.ListableURI, err error) {
			if err == nil && uc != nil {
				entry.SetText(uc.Path())
			}
		}, window).Show()
	})

	// Create a Table for File Display
	table := widget.NewTable(
		func() (int, int) { return len(files) + 1, 5 }, // Rows and columns

		func() fyne.CanvasObject {
			return widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: false})
		},

		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row == 0 { // Header row
				headers := []string{"Name", "Path", "Size (KB)", "Modified", "Hash"}
				label := cell.(*widget.Label)
				label.Wrapping = fyne.TextWrapWord // Enable word wrapping
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				row := files[id.Row-1]
				data := []string{row.Name, row.Path, row.Size, row.ModifiedTime, row.Hash}
				label.SetText(data[id.Col])
			}
		},
	)

	// Set Column Widths
	table.SetColumnWidth(0, 300) // Name
	table.SetColumnWidth(1, 350) // Path
	table.SetColumnWidth(2, 100) // Size
	table.SetColumnWidth(3, 220) // Modified Time
	table.SetColumnWidth(4, 300) // Hash

	// Wrap the table inside a scroll container and set its minimum size
	tableContainer := container.NewScroll(table)
	tableContainer.SetMinSize(fyne.NewSize(950, 400)) // Adjust table height

	// Scan button
	scanButton := widget.NewButton("Scan Directory", func() {
		dirPath := entry.Text
		// dirPath := "/home/iam5k/Documents"
		if dirPath == "" {
			output.Text = "Please select a directory first."
			output.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red for errors
			output.Refresh()
			return
		}

		spinner.Show() // Show loading spinner
		go func() {
			defer spinner.Hide()
			scannedFiles, err := filesystem.ScanDirectoryAndSaveMetadata(dirPath)
			if err != nil {
				output.Text = "Error: " + err.Error()
				output.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			} else {
				output.Text = "Directory scanned successfully!"
				output.Color = color.Black
				files = []FileData{}
				for _, f := range scannedFiles {
					files = append(files, FileData{
						Name:         f.Name,
						Path:         strings.TrimPrefix(f.Path, dirPath),
						Size:         strconv.FormatInt(f.Size/1024, 10) + " KB",
						ModifiedTime: f.ModifiedTime.Format(time.RFC1123),
						Hash:         f.Hash,
					})
				}
				table.Refresh() // Update UI
			}
			output.Refresh()
		}()
	})

	// Layout
	content := container.NewVBox(
		label,
		entry,
		selectButton,
		scanButton,
		output,
		spinner,
		tableContainer, // Use the wrapped scrollable table
	)

	window.SetContent(content)
	window.ShowAndRun()
}
