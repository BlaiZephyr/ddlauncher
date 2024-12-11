package main

import (
	"ddlauncher/backend"
	"ddlauncher/frontend"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow(backend.WindowTitle)

	backend.InitAppState()
	content := frontend.CreateMainUI(w)

	w.SetContent(content)
	w.Resize(frontend.DefaultWindowSize)
	w.ShowAndRun()
}
