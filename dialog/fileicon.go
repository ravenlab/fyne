package dialog

import (
	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/widget"
)

// NewFileIcon takes a filepath and creates an icon with an overlayed label using the detected mimetype and extension
// Deprecated: Use widget.NewFileIcon instead
func NewFileIcon(uri fyne.URI) *widget.FileIcon {
	return widget.NewFileIcon(uri)
}
