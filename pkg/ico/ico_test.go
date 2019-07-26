package ico

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
)

type badReader struct {
	blockRead bool
	data      []byte
	offset    int
}

func (r *badReader) Read(b []byte) (int, error) {
	if r.blockRead {
		return 0, errors.New("you cannot read from bad reader")
	}
	if r.offset == len(r.data) {
		return 0, io.EOF
	}
	n := copy(b, r.data[r.offset:])
	r.offset += n
	return n, nil
}

func (r *badReader) ReadAt(b []byte, offset int64) (int, error) {
	if r.blockRead {
		return 0, errors.New("you cannot read from bad reader")
	}
	if int(offset) >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(b, r.data[offset:])
	if n+r.offset == len(r.data) {
		return n, io.EOF
	}
	return n, nil
}

func TestDecodeAll(t *testing.T) {
	for _, tc := range []struct {
		Filename string
		Count    int
	}{
		{"icon.ico", 9},
	} {
		f, err := os.Open(filepath.Join("..", "..", "testdata", tc.Filename))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		g, err := DecodeAll(f)
		if err != nil {
			t.Fatal(err)
		}

		if len(g.entries) != tc.Count {
			t.Fatalf("wrong entries length; expected 9, got %d", len(g.entries))
		}
		if len(g.Images) != tc.Count {
			t.Fatalf("wrong images length; expected 9, got %d", len(g.Images))
		}
	}
}

func TestDecodeAll_invalidReader(t *testing.T) {
	for _, tc := range []Reader{
		&badReader{blockRead: true},
		&badReader{data: []byte{0, 0, 0, 0, 0, 0}},
		&badReader{data: []byte{0, 0, 1, 0, 0, 0}},
		&badReader{data: []byte{0, 0, 1, 0, 1, 0}},
		// &badReader{data: []byte{0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	} {
		if _, err := DecodeAll(tc); err == nil {
			t.Fatal("expected failure, got no error")
		}
	}
}
