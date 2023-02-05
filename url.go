package iokit

import (
	"github.com/pkgs-go/fu"
	"io"
	"strings"
)

type IoUrl struct {
	Url      string
	Schema   string
	Cache    Cache
	Observer AsyncUpload
	Metadata Metadata
}

type Metadata map[string]string
type AsyncUpload struct{ Notify func(url string, err error) }

func Url(url string, opts ...interface{}) StrictInputOutput {
	lurl := strings.ToLower(url)
	schema := ""
	if j := strings.Index(lurl, "://"); j > 0 {
		schema = lurl[:j]
	}
	return StrictInputOutput{IoUrl{
		url,
		schema,
		fu.Select(Cache(""), opts),
		fu.Select(AsyncUpload{nil}, opts),
		fu.Select(Metadata(nil), opts),
	}}
}

func (iourl IoUrl) Open() (rd io.ReadCloser, err error) {
	if iourl.Schema != "file" {
		return iourl.openUrlReader()
	}
	return File(iourl.Url[7:]).Open()
}

func (iourl IoUrl) Create() (hw Whole, err error) {
	if iourl.Schema != "file" {
		return iourl.createUrlWriter()
	}
	return File(iourl.Url[7:]).Create()
}
