package ico

import "io"

// Reader implements both io.Reader and io.ReaderAt.
type Reader interface {
	io.Reader
	io.ReaderAt
}
