package iokit

import (
	"fmt"
	"gotest.tools/assert"
	"io"
	"math/rand"
	"os"
	"testing"
)

func Test_Example(t *testing.T) {

	_ = os.Setenv("FILES", "/tmp/go-iokit-test-files")

	for _, url := range []string{
		//"s3://$do_test/test_example.txt",
		//"gs://$enctest/test_example.txt",
		"file://$files/test_example.txt"} {

		S := fmt.Sprintf(`Hello world! %d`, rand.Int())

		wh, err := Url(url).Create()
		assert.NilError(t, err)
		defer wh.End()
		_, err = wh.Write([]byte(S))
		assert.NilError(t, err)
		err = wh.Commit()
		assert.NilError(t, err)

		rd, err := Url(url).Open()
		assert.NilError(t, err)
		defer CloseNoError(rd)
		q, err := io.ReadAll(rd)
		assert.NilError(t, err)
		assert.Assert(t, string(q) == S)
	}
}
