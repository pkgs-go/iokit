package iokit

import (
	"archive/zip"
	"fmt"
	"io"
)

type zipfile struct {
	Arch     Input
	FileName string
}

func ZipFile(fileName string, arch Input) StrictInput {
	return StrictInput{zipfile{arch, fileName}}
}

func (q zipfile) Open() (f io.ReadCloser, err error) {
	var xf io.ReadCloser
	xf, err = q.Arch.Open()
	if err != nil {
		return
	}
	defer func() {
		if xf != nil {
			_ = xf.Close()
		}
	}()
	var r *zip.Reader
	if r, err = zip.NewReader(xf.(io.ReaderAt), FileSize(xf)); err != nil {
		return
	}
	for _, n := range r.File {
		if n.Name == q.FileName {
			zf, err := n.Open()
			if err != nil {
				return nil, err
			}
			xxf := xf
			xf = nil
			return Reader(zf, func() error {
				_ = zf.Close()
				return xxf.Close()
			}), nil
		}
	}
	return nil, fmt.Errorf("zip archive does not contain file " + q.FileName)
}

type zipout struct {
	Arch     Output
	FileName string
}

func Zip(fileName string, arch Output) StrictOutput {
	return StrictOutput{zipout{arch, fileName}}
}

type zipwhole struct {
	io.Writer
	zw *zip.Writer
	f  Whole
}

func (q zipout) Create() (w Whole, err error) {
	var xf Whole
	xf, err = q.Arch.Create()
	if err != nil {
		return
	}
	zw := zip.NewWriter(xf)
	fw, err := zw.Create(q.FileName)
	if err != nil {
		return
	}
	return &zipwhole{fw, zw, xf}, nil
}

func (z *zipwhole) End() {
	z.f.End()
}

func (z *zipwhole) Commit() error {
	CloseNoError(z.zw)
	return z.f.Commit()
}
