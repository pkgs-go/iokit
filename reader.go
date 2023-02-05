package iokit

import "io"

type reader struct {
	io.Reader
	close []func() error
}

func (r reader) Close() (err error) {
	for _, f := range r.close {
		if e := f(); e != nil {
			err = e
		}
	}
	return
}

func (r reader) Open() (io.ReadCloser, error) {
	return r, nil
}

func Reader(rd io.Reader, close ...func() error) io.ReadCloser {
	return reader{rd, close}
}
