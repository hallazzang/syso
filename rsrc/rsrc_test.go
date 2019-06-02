package rsrc

import (
	"bytes"
	"fmt"
	"io"
	"testing"
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
	data := newDummyBlob([]byte("hello"))

	s := New()
	s.rootDir.addSubdirectoryByName("Hola", 0xb8f).addDataEntryByName("Data", data)
	s.rootDir.addSubdirectoryByID(2, 0xc04).addDataEntryByID(3, data)

	b := new(bytes.Buffer)
	n, err := s.WriteTo(b)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(n)
	fmt.Printf("%q\n", b)
}
