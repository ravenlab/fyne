package widget_test

import (
	"testing"

	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/layout"
	"github.com/ravenlab/fyne/test"
	"github.com/ravenlab/fyne/theme"
	"github.com/ravenlab/fyne/widget"
)

func TestIcon_Layout(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, theme.LightTheme())

	for name, tt := range map[string]struct {
		resource fyne.Resource
	}{
		"empty": {},
		"resource": {
			resource: theme.CancelIcon(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			icon := &widget.Icon{
				Resource: tt.resource,
			}

			window := test.NewWindow(fyne.NewContainerWithLayout(layout.NewCenterLayout(), icon))
			window.Resize(icon.MinSize().Max(fyne.NewSize(150, 200)))

			test.AssertImageMatches(t, "icon/layout_"+name+".png", window.Canvas().Capture())

			window.Close()
		})
	}
}
