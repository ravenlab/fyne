package test

import (
	"os"

	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/storage"
)

type testStorage struct {
}

func (s *testStorage) RootURI() fyne.URI {
	return storage.NewURI("file://" + os.TempDir())
}
