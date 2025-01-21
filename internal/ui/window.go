package ui

import (
	"log"
	"strconv" // Import strconv for number-to-string conversion
	"suparna/internal/filesystem"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func Launch() {
	log.Println("Starting Suparna application...")

	// Use app.NewWithID to avoid Fyne preferences error
	myApp := app.NewWithID("suparna")
	window := myApp.NewWindow("Suparna File Manager")

	// Set the initial window size (increased to 1000x800 for better layout)
	window.Resize(fyne.NewSize(1000, 800))

	label := widget.NewLabel("Welcome to Suparna!")
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter directory path...")

	output := widget.NewMultiLineEntry()
	output.Disable() // Makes the MultiLineEntry read-only

	// Open folder selection dialog
	button := widget.NewButton("Select Directory", func() {
		fileDialog := dialog.NewFolderOpen(func(uc fyne.ListableURI, err error) {
			if err == nil && uc != nil {
				entry.SetText(uc.Path()) // Set the path of the selected folder to entry field
			}
		}, window)
		fileDialog.Show()
	})

	// Create table model to hold file data
	fileData := make([][]string, 0)
	fileTable := widget.NewTable(
		func() (int, int) { return len(fileData), 5 }, // Number of rows and columns
		func() fyne.CanvasObject {
			// Create headers for each column
			return container.NewHBox(
				widget.NewLabel("Name"),
				widget.NewLabel("Path"),
				widget.NewLabel("Size"),
				widget.NewLabel("Modified Time"),
				widget.NewLabel("Hash"),
			)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			// Update table cell contents dynamically
			label, ok := o.(*widget.Label) // Ensure we're working with a label
			if !ok {
				return
			}

			switch i.Col {
			case 0:
				label.SetText(fileData[i.Row][0]) // Set name
			case 1:
				label.SetText(fileData[i.Row][1]) // Set path
			case 2:
				label.SetText(fileData[i.Row][2]) // Set size
			case 3:
				label.SetText(fileData[i.Row][3]) // Set modified time
			case 4:
				label.SetText(fileData[i.Row][4]) // Set hash
			}
		},
	)

	// Button to scan directory
	scanButton := widget.NewButton("Scan Directory", func() {
		dirPath := entry.Text // Get the directory path from the entry field
		if dirPath == "" {
			output.SetText("Please select a directory first.")
			return
		}

		files, err := filesystem.ScanDirectoryAndSaveMetadata(dirPath)
		if err != nil {
			output.SetText("Error: " + err.Error())
			return
		}
		output.SetText("") // Clear the output before appending new text

		// Log files that were scanned and retrieved
		// log.Printf("Files to display in table: %+v", files)

		// Update the table with scanned files
		fileData = nil      // Reset previous data
		fileTable.Refresh() // Refresh the table after each file is added
		for _, file := range files {
			// log.Printf("File to display: %+v", file) //this displays correct data
			fileData = append(fileData, []string{
				file.Name,
				file.Path,
				strconv.FormatInt(file.Size, 10), // Correct conversion of int64 to string
				file.ModifiedTime.Format(time.RFC1123),
				file.Hash,
			})
		}

		// Log the file data array to ensure correct data
		log.Printf("FileData to display: %+v", fileData) //this displays correct data

		// Refresh table to reflect changes
		fileTable.Refresh() // Refresh the table after each file is added
	})
	log.Printf("FileTable to display: %+v", fileTable) //this displays incorrect data
	// Wrap the fileTable with a scroll container and define a height
	tableContainer := container.NewScroll(fileTable)
	tableContainer.SetMinSize(fyne.NewSize(0, 300)) // Set a minimum height for the table

	// Layout the UI components
	content := container.NewVBox(label, entry, button, scanButton, output, tableContainer)
	window.SetContent(content)
	window.ShowAndRun()
}
