package common

import (
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Blob represents arbitrary data that can be embedded in an object file.
type Blob interface {
	io.Reader
	Size() int64
}

type dataBlob struct {
	data   []byte
	offset int64
}

// NewBlob creates a blob from r by reading all data from it.
func NewBlob(r io.Reader) (Blob, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read data")
	}
	return &dataBlob{
		data: b,
	}, nil
}

func (b *dataBlob) Read(p []byte) (int, error) {
	n := copy(p[:], b.data[b.offset:])
	b.offset += int64(n)
	return n, nil
}

func (b *dataBlob) Size() int64 {
	return int64(len(b.data))
}
