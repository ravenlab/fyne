package dialog

import (
	"image/color"
	"testing"

	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/canvas"
	"github.com/ravenlab/fyne/test"
	"github.com/ravenlab/fyne/theme"
	"github.com/ravenlab/fyne/widget"
)

func TestShowCustom_ApplyTheme(t *testing.T) {
	test.ApplyTheme(t, theme.DarkTheme())
	w := test.NewWindow(canvas.NewRectangle(color.Transparent))
	w.Resize(fyne.NewSize(300, 200))

	label := widget.NewLabel("Content")
	label.Alignment = fyne.TextAlignCenter

	d := NewCustom("Title", "OK", label, w)

	d.Show()
	test.AssertImageMatches(t, "dialog-custom-dark.png", w.Canvas().Capture())

	test.ApplyTheme(t, theme.LightTheme())
	test.AssertImageMatches(t, "dialog-custom-light.png", w.Canvas().Capture())
}
