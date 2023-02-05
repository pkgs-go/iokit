package iokit

import (
	"bytes"
	"io"
	"io/ioutil"
)

type Input interface {
	Open() (io.ReadCloser, error)
}

type Output interface {
	Create() (Whole, error)
}

type InputOutput interface {
	Input
	Output
}

type StrictInputOutput struct{ InputOutput }

func (iox StrictInputOutput) MustCreate() StrictWhole {
	return StrictOutput{iox.InputOutput}.MustCreate()
}

func (iox StrictInputOutput) MustOpen() StrictReader {
	return StrictInput{iox.InputOutput}.MustOpen()
}

func (iox StrictInputOutput) ReadAll() ([]byte, error) {
	return StrictInput{iox.InputOutput}.ReadAll()
}

func (iox StrictInputOutput) MustReadAll() []byte {
	return StrictInput{iox.InputOutput}.MustReadAll()
}

func (iox StrictInputOutput) WriteAll(bs []byte) error {
	return StrictOutput{iox.InputOutput}.WriteAll(bs)
}

func (iox StrictInputOutput) MustWriteAll(bs []byte) {
	StrictOutput{iox.InputOutput}.MustWriteAll(bs)
}

type StrictOutput struct{ Output }

func (iox StrictOutput) MustCreate() StrictWhole {
	wr, err := iox.Create()
	if err != nil {
		panic(err)
	}
	return StrictWhole{wr}
}

func (iox StrictOutput) WriteAll(bs []byte) error {
	wr, err := iox.Create()
	if err != nil {
		return err
	}
	defer wr.End()
	_, err = io.Copy(wr, bytes.NewReader(bs))
	if err != nil {
		return err
	}
	return wr.Commit()
}

func (iox StrictOutput) MustWriteAll(bs []byte) {
	if err := iox.WriteAll(bs); err != nil {
		panic(err)
	}
}

type StrictWhole struct{ Whole }

func (lw StrictWhole) MustWrite(b []byte) {
	if _, err := lw.Write(b); err != nil {
		panic(err)
	}
}

func (lr StrictWhole) MustCommit() {
	if err := lr.Commit(); err != nil {
		panic(err)
	}
}

type StrictInput struct{ Input }

func (iox StrictInput) MustOpen() StrictReader {
	rd, err := iox.Open()
	if err != nil {
		panic(err)
	}
	return StrictReader{rd}
}

func (iox StrictInput) ReadAll() ([]byte, error) {
	rd, err := iox.Open()
	if err != nil {
		return nil, err
	}
	defer rd.Close()
	return ioutil.ReadAll(rd)
}

func (iox StrictInput) MustReadAll() []byte {
	bs, err := iox.ReadAll()
	if err != nil {
		panic(err)
	}
	return bs
}

type StrictReader struct{ io.ReadCloser }

func (lr StrictReader) MustReadAll() []byte {
	bs, err := ioutil.ReadAll(lr)
	if err != nil {
		panic(err)
	}
	return bs
}
