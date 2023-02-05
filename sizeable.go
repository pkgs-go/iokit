package iokit

import (
	"io"
	"os"
)

type Sizeable interface {
	Size() int64
}

func FileSize(rd io.Reader) int64 {
	if i, ok := rd.(Sizeable); ok {
		return i.Size()
	}
	if f, ok := rd.(*os.File); ok {
		st, err := f.Stat()
		if err != nil { // wine workaround
			st, err = os.Stat(f.Name())
		}
		if err == nil {
			return st.Size()
		}
	}
	return 0
}
