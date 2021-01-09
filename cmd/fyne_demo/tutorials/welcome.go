package tutorials

import (
	"net/url"

	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/canvas"
	"github.com/ravenlab/fyne/cmd/fyne_demo/data"
	"github.com/ravenlab/fyne/container"
	"github.com/ravenlab/fyne/widget"
)

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func welcomeScreen(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneScene)
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(171, 125))
	} else {
		logo.SetMinSize(fyne.NewSize(228, 167))
	}

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Welcome to the Fyne toolkit demo app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		container.NewHBox(
			widget.NewHyperlink("github.com/ravenlab", parseURL("https://github.com/ravenlab/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("documentation", parseURL("https://github.com/ravenlab/develop/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("sponsor", parseURL("https://github.com/sponsors/github.com/ravenlab")),
		),
	))
}
