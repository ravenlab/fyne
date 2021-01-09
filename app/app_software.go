// +build ci nacl

package app

import (
	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/internal/painter/software"
	"github.com/ravenlab/fyne/test"
)

// NewWithID returns a new app instance using the test (headless) driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
	return newAppWithDriver(test.NewDriverWithPainter(software.NewPainter()), id)
}
