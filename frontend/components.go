package frontend

import (
	"ddlauncher/backend"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createMainTabs() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")),
		container.NewTabItem("AntiBot", widget.NewLabel("AntiBot")),
		container.NewTabItem("Translator", widget.NewLabel("Translator")),
	)
	return tabs
}

func createGameStarterButton(w fyne.Window) *widget.Button {
	return widget.NewButton(backend.GameTitle, func() {
		err := backend.RunGameCommand()
		if err != nil {
			fmt.Println(err)
		}
	})
}

func createVersionControlButton(w fyne.Window) fyne.CanvasObject {
	tags, err := backend.FetchGitHubTags()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	versionSelect := widget.NewSelect(tags, func(selected string) {
		backend.State.CurrentVersion = selected
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
