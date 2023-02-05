package iokit

import (
	"fmt"
	"gotest.tools/assert"
	"math/rand"
	"testing"
)

func Test_Zip(t *testing.T) {
	S := fmt.Sprintf("test string %v", rand.Int())
	func() {
		w := Zip("test.txt", Cache("test.zip").File()).MustCreate()
		defer w.End()
		w.MustWrite([]byte(S))
		w.MustCommit()
	}()
	s := func() string {
		r := ZipFile("test.txt", Cache("test.zip").File()).MustOpen()
		defer CloseNoError(r)
		return string(r.MustReadAll())
	}()
	assert.Assert(t, s == S)
}

func Test_Zip_Fail(t *testing.T) {
	S := fmt.Sprintf("test string %v", rand.Int())
	func() {
		w := Zip("test.txt", Cache("test.zip").File()).MustCreate()
		defer w.End()
		w.MustWrite([]byte(S))
		// no commit here
	}()
	_, err := ZipFile("test.txt", Cache("test.zip").File()).Open()
	assert.Assert(t, err != nil)
}
