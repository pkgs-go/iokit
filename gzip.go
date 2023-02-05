package iokit

import (
	"compress/gzip"
	"io"
)

type gzipf struct {
	Output
}

func Gzip(output Output) StrictOutput {
	return StrictOutput{gzipf{output}}
}

type gzipwr struct {
	wl io.WriteCloser
	wh Whole
}

func (gz gzipf) Create() (wh Whole, err error) {
	w, err := gz.Output.Create()
	if err != nil {
		return
	}
	wx := gzip.NewWriter(w)
	return &gzipwr{wx, w}, nil
}

func (gz *gzipwr) Commit() error {
	if err := gz.wl.Close(); err != nil {
		return err
	}
	return gz.wh.Commit()
}

func (gz *gzipwr) End() {
	_ = gz.wl.Close()
	gz.wh.End()
}

func (gz *gzipwr) Write(b []byte) (int, error) {
	return gz.wl.Write(b)
}
