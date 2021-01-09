// +build !ci,!nacl,!android,!ios,!mobile

package app

import (
	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/internal/driver/glfw"
)

// NewWithID returns a new app instance using the appropriate runtime driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
	return newAppWithDriver(glfw.NewGLDriver(), id)
}
