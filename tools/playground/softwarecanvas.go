package playground

import (
	"github.com/ravenlab/fyne/internal/painter/software"
	"github.com/ravenlab/fyne/test"
)

// NewSoftwareCanvas creates a new canvas in memory that can render without hardware support
func NewSoftwareCanvas() test.WindowlessCanvas {
	return test.NewCanvasWithPainter(software.NewPainter())
}
