// Package fyne describes the objects and components available to any Fyne app.
// These can all be created, manipulated and tested without rendering (for speed).
// Your main package should use the app package to create an application with
// a default driver that will render your UI.
//
// A simple application may look like this:
//
//   package main
//
//   import "github.com/ravenlab/fyne/app"
//   import "github.com/ravenlab/fyne/widget"
//
//   func main() {
//   	a := app.New()
//   	w := a.NewWindow("Hello")
//
//   	hello := widget.NewLabel("Hello Fyne!")
//   	w.SetContent(widget.NewVBox(
//   		hello,
//   		widget.NewButton("Hi!", func() {
//   			hello.SetText("Welcome :)")
//   		}),
//   	))
//
//   	w.ShowAndRun()
//   }
package fyne // import "github.com/ravenlab/fyne"
