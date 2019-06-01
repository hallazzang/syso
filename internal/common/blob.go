package common

import "io"

// Blob represents arbitary data that can be embedded in an object file.
type Blob interface {
	io.Reader
	Size() int64
}

func foo() {

}
