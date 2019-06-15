package syso

import (
	"fmt"
	"io"
	"os"
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
	b := newDummyBlob([]byte("helloworld"))
	if err := r.AddIconByID(1, b); err != nil {
		t.Fatal(err)
	}

	if err := c.AddSection(r); err != nil {
		t.Fatal(err)
	}

	f, err := os.Create("out.obj")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	n, err := c.WriteTo(f)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(n)
}
