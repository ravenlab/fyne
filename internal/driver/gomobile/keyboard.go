package gomobile

import (
	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/driver/mobile"
	"github.com/github.com/ravenlab/mobile/app"
)

func showVirtualKeyboard(keyboard mobile.KeyboardType) {
	if driver, ok := fyne.CurrentApp().Driver().(*mobileDriver); ok {
		driver.app.ShowVirtualKeyboard(app.KeyboardType(keyboard))
	}
}

func hideVirtualKeyboard() {
	if driver, ok := fyne.CurrentApp().Driver().(*mobileDriver); ok {
		driver.app.HideVirtualKeyboard()
	}
}
