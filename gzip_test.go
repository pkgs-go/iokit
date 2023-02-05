package iokit

import (
	"fmt"
	"gotest.tools/assert"
	"math/rand"
	"testing"
)

func Test_Gzip(t *testing.T) {
	S := fmt.Sprintf("test string %v", rand.Int())
	func() {
		w := Gzip(Cache("test.gz").File()).MustCreate()
		defer w.End()
		w.MustWrite([]byte(S))
		w.MustCommit()
	}()
	s := func() string {
		r := Compressed(Cache("test.gz").File()).MustOpen()
		defer CloseNoError(r)
		return string(r.MustReadAll())
	}()
	assert.Assert(t, s == S)
}

func Test_Gzip_Fail(t *testing.T) {
	S := fmt.Sprintf("test string %v", rand.Int())
	func() {
		w := Gzip(Cache("test.gz").File()).MustCreate()
		defer w.End()
		w.MustWrite([]byte(S))
		// no commit here
	}()
	_, err := Compressed(Cache("test.gz").File()).Open()
	assert.Assert(t, err != nil)
}
