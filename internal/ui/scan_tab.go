package ui

import (
	"image/color"
	"log"
	"suparna/internal/database"
	"suparna/internal/filesystem"
	"sync"

	"fyne.io/fyne/v2"
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
var scanningMutex sync.Mutex
var scanningAborted bool

// CreateScanTab defines the UI for the "Scan and Add Directory" tab.
func CreateScanTab() fyne.CanvasObject {
	log.Println("Initializing Scan Tab...")

	// Components
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter directory path...")

	output := canvas.NewText("", color.Black)
	progress := widget.NewProgressBar() // Linear progress bar
	progress.Hide()

	currentFileLabel := widget.NewLabel("") // Label for current file being scanned
	currentFileLabel.Hide()

	// Stop scan flag
	isScanning := false

	// Directory selection button
	selectButton := widget.NewButton("Select Directory", func() {
		dialog.NewFolderOpen(func(uc fyne.ListableURI, err error) {
			if err == nil && uc != nil {
				entry.SetText(uc.Path())
			}
		}, fyne.CurrentApp().Driver().AllWindows()[0]).Show()
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
				label.SetText(headers[id.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				row := files[id.Row-1]
				data := []string{row.Name, row.Path, row.Size, row.ModifiedTime, row.Hash}
				label.Wrapping = fyne.TextWrapWord
				label.SetText(data[id.Col])
			}
		},
	)
	// Scan button
	scanButton := widget.NewButton("Scan Directory", nil) // Initialize without action
	scanButton.OnTapped = func() {
		dirPath := entry.Text
		if dirPath == "" {
			output.Text = "Please select a directory first."
			output.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red for errors
			output.Refresh()
			return
		}

		if isScanning { // Stop scanning logic
			scanningMutex.Lock()
			scanningAborted = true
			scanningMutex.Unlock()
			isScanning = false
			scanButton.SetText("Scan Directory")
			scanButton.Importance = widget.HighImportance
			output.Text = "Scanning stopped by user."
			output.Refresh()
			return
		}

		// Start scanning logic
		isScanning = true
		scanningAborted = false
		scanButton.SetText("Stop Scan")
		scanButton.Importance = widget.DangerImportance // Red button

		progress.Show()
		currentFileLabel.Show()
		output.Text = ""
		output.Refresh()

		go func() {
			files = []FileData{} // Clear previous data

			err := filesystem.ScanDirectory(dirPath, func(currentFile string, progressValue float64) {
				scanningMutex.Lock()
				if scanningAborted {
					scanningMutex.Unlock()
					return
				}
				scanningMutex.Unlock()

				// Update UI dynamically
				currentFileLabel.SetText("Scanning: " + currentFile)
				progress.SetValue(progressValue)
				currentFileLabel.Refresh()
				progress.Refresh()
			})

			// ✅ Handle error properly
			if err != nil {
				log.Printf("Error during scan: %v", err)              // Log to console
				output.Text = "Error: " + err.Error()                 // Show error in UI
				output.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red text
				output.Refresh()
				return
			}

			// ✅ Fetch newly scanned files from the database
			db := database.GetDB()
			rows, err := db.Query("SELECT name, path, size, modified_time, hash FROM files ORDER BY modified_time DESC")
			if err != nil {
				log.Printf("Database query error: %v", err)
				output.Text = "DB Error: " + err.Error()
				output.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
				output.Refresh()
				return
			}
			defer rows.Close()

			// Read rows into files slice
			for rows.Next() {
				var f FileData
				var modifiedTimeStr string
				err = rows.Scan(&f.Name, &f.Path, &f.Size, &modifiedTimeStr, &f.Hash)
				if err != nil {
					log.Printf("Error scanning row: %v", err)
					continue
				}
				f.ModifiedTime = modifiedTimeStr
				files = append(files, f)
			}

			// Ensure table updates on UI thread
			fyne.CurrentApp().Driver().AllWindows()[0].Canvas().Refresh(table)
			table.Refresh()
		}()

	}

	// Explicitly set column widths
	table.SetColumnWidth(0, 300) // Name
	table.SetColumnWidth(1, 350) // Path
	table.SetColumnWidth(2, 100) // Size
	table.SetColumnWidth(3, 220) // Modified Time
	table.SetColumnWidth(4, 300) // Hash

	// Ensure the table has proper dimensions
	tableContainer := container.NewScroll(table)
	tableContainer.SetMinSize(fyne.NewSize(950, 400)) // Adjust table height and width

	// Layout the UI
	buttonRow := container.NewGridWithColumns(2, selectButton, scanButton) // Inline buttons
	content := container.NewVBox(
		widget.NewLabel("Scan and Add Directory"),
		entry,
		buttonRow,
		progress,
		currentFileLabel,
		output,
		tableContainer, // Table inside a scrollable container
	)

	return content
}
