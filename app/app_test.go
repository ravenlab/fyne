package app

import (
	"net/url"
	"os/exec"
	"strings"
	"testing"

	"github.com/ravenlab/fyne"
	_ "github.com/ravenlab/fyne/test"
	"github.com/stretchr/testify/assert"
)

func TestDummyApp(t *testing.T) {
	app := NewWithID("io.fyne.test")

	app.Quit()
}

func TestCurrentApp(t *testing.T) {
	app := NewWithID("io.fyne.test")

	assert.Equal(t, app, fyne.CurrentApp())
}

func TestFyneApp_UniqueID(t *testing.T) {
	appID := "io.fyne.test"
	app := NewWithID(appID)

	assert.Equal(t, appID, app.UniqueID())
}

func TestFyneApp_OpenURL(t *testing.T) {
	opened := ""
	app := NewWithID("io.fyne.test")
	app.(*fyneApp).exec = func(cmd string, arg ...string) *exec.Cmd {
		opened = arg[len(arg)-1]
		return exec.Command("")
	}

	urlStr := "https://github.com/ravenlab"
	u, _ := url.Parse(urlStr)
	err := app.OpenURL(u)

	if err != nil && strings.Contains(err.Error(), "unknown operating system") {
		return // when running in CI mode we don't actually open URLs...
	}

	assert.Equal(t, urlStr, opened)
}
