package iokit

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"github.com/pkgs-go/fu"
	"github.com/ulikunitz/xz"
	"io"
	"os"
)

type compressed struct{ input interface{} }

func Compressed(arch interface{}) StrictInput {
	return StrictInput{compressed{arch}}
}

func (q compressed) Open() (f io.ReadCloser, err error) {
	var xf io.ReadCloser
	if e, ok := q.input.(Input); ok {
		xf, err = e.Open()
	} else {
		xf, err = os.Open(q.input.(string))
	}
	if err != nil {
		return
	}
	dc := decompress(xf)
	return Reader(dc.Run(),
		func() error {
			e := dc.Close()
			err := xf.Close()
			if e != nil {
				return e
			}
			return err
		}), nil
}

const decompressorBufferSize = 32 * 1024

type errReader struct{ err error }

func (errReader) Read(p []byte) (n int, err error) {
	return 0, err
}

func (errReader) Close() error {
	return nil
}

func Decompress(source interface{}) io.ReadCloser {
	if q, ok := source.([]byte); ok {
		return decompress(bytes.NewBuffer(q)).Run()
	}
	if q, ok := source.(io.Reader); ok {
		return decompress(q).Run()
	}
	return &errReader{io.ErrUnexpectedEOF}
}

const bufferCounts = 4

type decomp struct {
	reader   io.Reader
	canclose bool
	buffer   [bufferCounts][decompressorBufferSize]byte
	size     [bufferCounts]int
	side     int
	count    int
	err      error
	stop     chan struct{}
	next     chan int
}

func (d *decomp) Read(p []byte) (n int, err error) {
	for d.count >= d.size[d.side] {
		var ok bool
		d.side, ok = <-d.next
		if !ok {
			return 0, d.err
		}
		d.count = 0
	}
	k := fu.Min(len(p), d.size[d.side]-d.count)
	copy(p[:k], d.buffer[d.side][d.count:d.count+k])
	d.count += k
	return k, nil
}

func decompressor(rd io.Reader, canclose bool) *decomp {
	return &decomp{
		reader:   rd,
		canclose: canclose,
		stop:     make(chan struct{}),
		next:     make(chan int, bufferCounts-2),
		err:      io.EOF,
	}
}

func (d *decomp) Run() io.ReadCloser {
	go func() {
		stop := d.stop
		side := 1
		for {
			k, err := io.ReadFull(d.reader, d.buffer[side][:])
			if k != 0 {
				d.size[side] = k
				select {
				case d.next <- side:
					side = (side + 1) % bufferCounts
				case <-stop:
					return
				}
			} else {
				d.err = err
				close(d.next)
				return
			}
		}
	}()
	return d
}

func (d *decomp) Close() error {
	if d.stop != nil {
		close(d.stop)
	}
	if d.canclose {
		if c, ok := d.reader.(io.Closer); ok {
			_ = c.Close()
		}
	}
	return nil
}

func decompress(rd io.Reader) *decomp {
	qr := bufio.NewReaderSize(rd, 32*1024)
	if b, err := qr.Peek(4); err != nil {
		return decompressor(&errReader{err}, false)
	} else {
		// BZIP2
		if b[0] == 0x42 && b[1] == 0x5A && b[2] == 0x68 {
			return decompressor(bzip2.NewReader(rd), false)
		}
		// GZIP
		if b[0] == 0x1f && b[1] == 0x8b {
			r, err := gzip.NewReader(qr)
			if err != nil {
				return decompressor(&errReader{err}, false)
			}
			return decompressor(r, true)
		}
		// XZ/LZMA2
		if b[0] == 0xFD && b[1] == 0x37 && b[2] == 0x7A && b[3] == 0x58 {
			r, err := xz.NewReader(qr)
			if err != nil {
				return decompressor(&errReader{err}, false)
			}
			return decompressor(r, true)
		}
		return decompressor(qr, false)
	}
}
