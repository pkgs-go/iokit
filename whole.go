package iokit

import "io"

type Whole interface {
	io.Writer
	Commit() error
	End()
}

type whole struct{ io.Writer }
type Fallible interface {
	Fail()
}

func (t *whole) End() {
	if t.Writer != nil {
		if f, ok := t.Writer.(Fallible); ok {
			f.Fail()
		} else if c, ok := t.Writer.(io.Closer); ok {
			_ = c.Close()
		}
		t.Writer = nil
	}
}

func (t *whole) Commit() (err error) {
	if c, ok := t.Writer.(io.Closer); ok {
		err = c.Close()
	}
	t.Writer = nil
	return
}
