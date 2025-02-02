package dialog

import (
	"testing"

	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/test"
	"github.com/ravenlab/fyne/theme"
)

func Test_colorWheel_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	wheel := newColorWheel(nil)
	wheel.SetHSLA(180, 100, 50, 255)
	window := test.NewWindow(wheel)
	window.Resize(wheel.MinSize().Max(fyne.NewSize(100, 100)))

	test.AssertImageMatches(t, "color/wheel_layout.png", window.Canvas().Capture())

	window.Close()
}
