package common

import (
	"bytes"
	"io"
	"testing"

	"github.com/pkg/errors"
)

type badReader struct{}

func (*badReader) Read([]byte) (int, error) {
	return 0, errors.New("you cannot read from bad reader")
}

func TestNewBlob(t *testing.T) {
	for _, tc := range []struct {
		Reader io.Reader
		Size   int64
		Data   []byte
	}{
		{bytes.NewReader([]byte{}), 0, []byte{}},
		{bytes.NewReader([]byte("\x00\x00\x00\x00")), 4, []byte{0, 0, 0, 0}},
	} {
		b, err := NewBlob(tc.Reader)
		if err != nil {
			t.Fatal(err)
		}
		if size := b.Size(); size != tc.Size {
			t.Fatalf("wrong size; expected %d, got %d", tc.Size, size)
		}
		buf := make([]byte, tc.Size)
		n, err := b.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		if n != int(tc.Size) {
			t.Fatalf("wrong read length; expected %d, got %d", tc.Size, n)
		}
		if !bytes.Equal(buf, tc.Data) {
			t.Fatalf("wrong data; expected %+q, got %+q", tc.Data, buf)
		}
	}
}

func TestNewBlob_invalidReader(t *testing.T) {
	for _, tc := range []io.Reader{
		&badReader{},
	} {
		b, err := NewBlob(tc)
		if err == nil {
			t.Fatalf("expected failure for reader type %T, got no error", b)
		}
		if b != nil {
			t.Fatal("must not create dataBlob instance")
		}
	}
}
