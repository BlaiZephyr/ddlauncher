package frontend

import (
	"ddlauncher/backend"
	"io"
	"log"
	"os"
	"sync"

	dialog2 "fyne.io/fyne/v2/dialog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ConsoleWriter struct {
	entry *widget.Entry
	mu    sync.Mutex
}

func (cw *ConsoleWriter) Write(p []byte) (n int, err error) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	cw.entry.SetText(cw.entry.Text + string(p))

	cw.entry.CursorRow = len(cw.entry.Text)

	return len(p), nil
}

func CreateConsoleOutput() (*widget.Entry, *ConsoleWriter) {
	consoleEntry := widget.NewMultiLineEntry()
	consoleEntry.Disable() // Make it read-only
	consoleEntry.SetMinRowsVisible(10)

	consoleWriter := &ConsoleWriter{entry: consoleEntry}
	return consoleEntry, consoleWriter
}

func RedirectStdoutAndStderr(consoleWriter *ConsoleWriter) {
	oldStdout := os.Stdout
	read, write, _ := os.Pipe()
	os.Stdout = write

	oldStderr := os.Stderr
	errR, errW, _ := os.Pipe()
	os.Stderr = errW

	go func() {
		io.Copy(consoleWriter, read)
		read.Close()
		os.Stdout = oldStdout
	}()

	go func() {
		io.Copy(consoleWriter, errR)
		errR.Close()
		os.Stderr = oldStderr
	}()
}

func createMainTabs(consoleOutput *widget.Entry) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")),
		container.NewTabItem("Console", consoleOutput),
	)
	return tabs
}

func createGameStarterButton(w fyne.Window) *widget.Button {
	return widget.NewButton(backend.GameTitle, func() {
		err := backend.RunGameCommand()
		if err != nil {
			dialog := dialog2.NewError(err, w)
			dialog.Show()
		}
	})
}

func createVersionControlButton() fyne.CanvasObject {
	versionSelect := widget.NewSelect([]string{"Loading..."}, nil)
	versionSelect.Disable()

	backend.FetchGitHubTagsAsync(func(tags []string, err error) {
		if err != nil {
			log.Println("Failed to fetch tags:", err)
			versionSelect.SetOptions([]string{"Error fetching versions"})
			return
		}

		versionSelect.SetOptions(tags)
		versionSelect.Enable()

		if len(tags) > 0 && backend.State.CurrentVersion != tags[0] {
			backend.State.CurrentVersion = tags[0]
			versionSelect.SetSelected(tags[0])
		}

		versionSelect.OnChanged = func(selected string) {
			backend.State.CurrentVersion = selected
		}
	})

	scrollContainer := container.NewVScroll(versionSelect)
	scrollContainer.SetMinSize(fyne.NewSize(200, 40))
	label := widget.NewLabel("Select Game Version")
	content := container.NewVBox(label, scrollContainer)
	return content
}

func createGameImage() fyne.CanvasObject {
	imagePath := "./images/tee.png"
	image := canvas.NewImageFromFile(imagePath)
	image.FillMode = canvas.ImageFillOriginal

	return container.New(layout.NewGridLayout(1), image)
}

func CreateMainContent(w fyne.Window) fyne.CanvasObject {
	consoleOutput, consoleWriter := CreateConsoleOutput()

	RedirectStdoutAndStderr(consoleWriter)

	mainTabs := createMainTabs(consoleOutput)

	gameStarterButton := createGameStarterButton(w)
	versionControlButton := createVersionControlButton()
	gameImage := createGameImage()

	content := container.NewVBox(
		gameImage,
		gameStarterButton,
		versionControlButton,
		mainTabs,
	)

	return content
}
