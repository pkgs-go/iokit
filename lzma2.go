package iokit

import (
	"github.com/ulikunitz/xz"
	"io"
)

type lzma2 struct {
	Output
}

func Lzma2(output Output) StrictOutput {
	return StrictOutput{lzma2{output}}
}

type lzma2wr struct {
	wl io.WriteCloser
	wh Whole
}

func (lz lzma2) Create() (wh Whole, err error) {
	w, err := lz.Output.Create()
	if err != nil {
		return
	}
	wx, err := xz.NewWriter(w)
	if err != nil {
		return
	}
	return &lzma2wr{wx, w}, nil
}

func (lz *lzma2wr) Commit() error {
	if err := lz.wl.Close(); err != nil {
		return err
	}
	return lz.wh.Commit()
}

func (lz *lzma2wr) End() {
	_ = lz.wl.Close()
	lz.wh.End()
}

func (lz *lzma2wr) Write(b []byte) (int, error) {
	return lz.wl.Write(b)
}
