package iokit

import (
	"fmt"
	"io"
)

type Resettable interface {
	Reset() error
}

func ResetFile(rd io.Reader) error {
	if i, ok := rd.(Resettable); ok {
		return i.Reset()
	}
	if i, ok := rd.(io.Seeker); ok {
		_, err := i.Seek(0, 0)
		return err
	}
	return fmt.Errorf("file is not resettable")
}
