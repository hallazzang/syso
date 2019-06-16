package syso

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/hallazzang/syso/pkg/coff"
	"github.com/hallazzang/syso/pkg/ico"
	"github.com/hallazzang/syso/pkg/rsrc"
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
	f, err := os.Open(filepath.Join("testdata", "icon.ico"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	g, err := ico.DecodeAll(f)
	if err != nil {
		t.Fatal(err)
	}

	for i, img := range g.Images {
		img.ID = i + 100
	}

	r := rsrc.New()
	if err := r.AddIconsByID(1, g); err != nil {
		t.Fatal(err)
	}

	c := coff.New()
	if err := c.AddSection(r); err != nil {
		t.Fatal(err)
	}

	of, err := os.Create(filepath.Join("testdata", "syso.syso"))
	if err != nil {
		t.Fatal(err)
	}
	defer of.Close()

	if _, err := c.WriteTo(of); err != nil {
		t.Fatal(err)
	}
}
