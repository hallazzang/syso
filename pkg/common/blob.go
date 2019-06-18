package common

import (
	"io"
	"os"
)

// Blob represents arbitrary data that can be embedded in an object file.
type Blob interface {
	io.Reader
	Size() int64
}

// FileBlob wraps os.File and implements Blob interface.
type FileBlob struct {
	*os.File
}

// NewFileBlob returns new FileBlob that wraps file f.
func NewFileBlob(f *os.File) *FileBlob {
	return &FileBlob{
		File: f,
	}
}

// Size returns file's size.
func (b *FileBlob) Size() int64 {
	fs, err := b.File.Stat()
	if err != nil {
		panic(err)
	}
	return fs.Size()
}
