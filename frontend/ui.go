package frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var DefaultWindowSize = fyne.NewSize(500, 500)

func CreateMainUI(w fyne.Window) fyne.CanvasObject {

	vBox := container.New(layout.NewVBoxLayout())

	vBox.Add(createGameImage())
	vBox.Add(createMainTabs())
	vBox.Add(createGameStarterButton(w))
	vBox.Add(createVersionControlButton(w))

	return vBox
}
