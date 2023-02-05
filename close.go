package iokit

import (
	"github.com/pkgs-go/error2"
	"io"
)

func CloseNoError(f io.Closer) {
	_ = f.Close() // ignore error
}

func MustClose(f io.Closer) {
	if err := f.Close(); err != nil {
		panic(error2.With(err, "an unexpected error occur during closing file"))
	}
}
