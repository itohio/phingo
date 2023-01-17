package app

import (
	"fmt"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func (a *App) MakeURI(base string) (fyne.URI, error) {
	return storage.Child(a.app.Storage().RootURI(), base)
}

func (a *App) Reader(base string) (io.ReadCloser, error) {
	uri, err := a.MakeURI(base)
	if err != nil {
		return nil, err
	}

	if ok, err := storage.CanRead(uri); !ok || err != nil {
		return nil, fmt.Errorf("Cannot read: %v", err)
	}

	reader, err := storage.Reader(uri)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func (a *App) Writer(base string) (io.WriteCloser, error) {
	uri, err := a.MakeURI(base)
	if err != nil {
		return nil, err
	}

	if ok, err := storage.CanWrite(uri); !ok || err != nil {
		return nil, fmt.Errorf("Cannot write: %v", err)
	}

	writer, err := storage.Writer(uri)
	if err != nil {
		return nil, err
	}

	return writer, nil
}
