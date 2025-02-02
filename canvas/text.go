package canvas

import (
	"image/color"

	"github.com/ravenlab/fyne"
)

// Declare conformity with CanvasObject interface
var _ fyne.CanvasObject = (*Text)(nil)

// Text describes a text primitive in a Fyne canvas.
// A text object can have a style set which will apply to the whole string.
// No formatting or text parsing will be performed
type Text struct {
	baseObject
	Alignment fyne.TextAlign // The alignment of the text content

	Color     color.Color    // The main text draw color
	Text      string         // The string content of this Text
	TextSize  int            // Size of the text - if the Canvas scale is 1.0 this will be equivalent to point size
	TextStyle fyne.TextStyle // The style of the text content
}

// MinSize returns the minimum size of this text object based on its font size and content.
// This is normally determined by the render implementation.
func (t *Text) MinSize() fyne.Size {
	return fyne.MeasureText(t.Text, t.TextSize, t.TextStyle)
}

// SetMinSize has no effect as the smallest size this canvas object can be is based on its font size and content.
func (t *Text) SetMinSize(size fyne.Size) {
	// no-op
}

// Refresh causes this object to be redrawn in it's current state
func (t *Text) Refresh() {
	Refresh(t)
}

// NewText returns a new Text implementation
func NewText(text string, color color.Color) *Text {
	return &Text{
		Color:    color,
		Text:     text,
		TextSize: fyne.CurrentApp().Settings().Theme().TextSize(),
	}
}
