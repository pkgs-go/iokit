package iokit

import "io"

type writer struct {
	io.Writer
	close []func(bool) error
}

func Writer(wr io.Writer, close ...func(bool) error) writer {
	return writer{wr, close}
}

func (w writer) Create() (Whole, error) {
	// yep, the ptr to copy of call receiver
	return &w, nil
}

func (w *writer) End() {
	for _, f := range w.close {
		_ = f(false)
	}
	return
}

func (w *writer) Commit() (err error) {
	for _, f := range w.close {
		if e := f(true); e != nil {
			err = e
		}
	}
	w.close = nil
	return
}
