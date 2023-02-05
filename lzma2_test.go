package iokit

import (
	"fmt"
	"gotest.tools/assert"
	"math/rand"
	"testing"
)

func Test_Lzma2(t *testing.T) {
	S := fmt.Sprintf("test string %v", rand.Int())
	func() {
		w := Lzma2(Cache("test.lzma2").File()).MustCreate()
		defer w.End()
		w.MustWrite([]byte(S))
		w.MustCommit()
	}()
	s := func() string {
		r := Compressed(Cache("test.lzma2").File()).MustOpen()
		defer r.Close()
		return string(r.MustReadAll())
	}()
	assert.Assert(t, s == S)
}
