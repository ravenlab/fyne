package widget

import (
	"testing"

	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/test"
	"github.com/ravenlab/fyne/theme"

	"github.com/stretchr/testify/assert"
)

func TestRadioItem_FocusIndicator_Centered_Vertically(t *testing.T) {
	item := newRadioItem("Hello", nil)
	render := test.WidgetRenderer(item).(*radioItemRenderer)
	render.Layout(fyne.NewSize(200, 100))

	focusIndicatorSize := theme.IconInlineSize() + theme.Padding()*2
	heightCenterOffset := (100 - focusIndicatorSize) / 2
	assert.Equal(t, fyne.NewPos(0, heightCenterOffset), render.focusIndicator.Position1)
}
