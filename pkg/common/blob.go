package common

import "io"

// Blob represents arbitrary data that can be embedded in an object file.
type Blob interface {
	io.Reader
	Size() int64
}
