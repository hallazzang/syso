package syso

import (
	"io"
	"testing"

	"github.com/hallazzang/syso/coff"
	"github.com/hallazzang/syso/rsrc"
)

type dummyBlob struct {
	data []byte
	size int
}

func newDummyBlob(data []byte) *dummyBlob {
	return &dummyBlob{
		data: data,
		size: len(data),
	}
}

func (r *dummyBlob) Read(b []byte) (int, error) {
	copy(b[:], r.data[:])
	return r.size, io.EOF
}

func (r *dummyBlob) Size() int64 {
	return int64(r.size)
}

func TestBasic(t *testing.T) {
	c := coff.New()

	r := rsrc.New()
	if err := c.AddSection(r); err != nil {
		t.Fatal(err)
	}

	b := newDummyBlob([]byte("hello"))
	if err := r.AddIconByID(1, b); err != nil {
		t.Fatal(err)
	}
}
