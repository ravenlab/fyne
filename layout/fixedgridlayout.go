package layout

import (
	"github.com/ravenlab/fyne"
)

// NewFixedGridLayout returns a new FixedGridLayout instance
//
// Deprecated: use the replacement NewGridWrapLayout. This method will be removed in 2.0.
func NewFixedGridLayout(size fyne.Size) fyne.Layout {
	return NewGridWrapLayout(size)
}
