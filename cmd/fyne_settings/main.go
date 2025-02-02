package main

import (
	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/app"
	"github.com/ravenlab/fyne/cmd/fyne_settings/settings"
	"github.com/ravenlab/fyne/widget"
)

func main() {
	s := settings.NewSettings()

	a := app.New()
	w := a.NewWindow("Fyne Settings")

	appearance := s.LoadAppearanceScreen(w)
	tabs := widget.NewTabContainer(
		&widget.TabItem{Text: "Appearance", Icon: s.AppearanceIcon(), Content: appearance})
	tabs.SetTabLocation(widget.TabLocationLeading)
	w.SetContent(tabs)

	w.Resize(fyne.NewSize(480, 480))
	w.ShowAndRun()
}
