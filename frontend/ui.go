package frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var DefaultWindowSize = fyne.NewSize(1000, 1000)

func CreateMainUI(w fyne.Window) fyne.CanvasObject {

	vBox := container.New(layout.NewVBoxLayout())

	vBox.Add(CreateMainContent(w))

	return vBox
}
