package glfw

import "github.com/ravenlab/fyne"

func (w *window) platformResize(canvasSize fyne.Size) {
	w.canvas.Resize(canvasSize)
	w.canvas.Refresh(w.canvas.content)
}
