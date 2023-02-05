package iokit

import "io"

type CloserChain []io.Closer

func (c CloserChain) Close() error {
	for _, x := range c {
		if x != nil {
			_ = x.Close()
		}
	}
	return nil
}
