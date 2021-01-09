// +build !ios,!android

package gomobile

import "github.com/ravenlab/fyne"

func canListURI(fyne.URI) bool {
	return false
}

func listURI(fyne.URI) ([]fyne.URI, error) {
	return nil, nil
}
