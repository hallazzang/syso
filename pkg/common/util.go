package common

import (
	"encoding/binary"
	"io"
)

// BinaryWriteTo writes v to w in little endian format.
func BinaryWriteTo(w io.Writer, v interface{}) (int64, error) {
	if err := binary.Write(w, binary.LittleEndian, v); err != nil {
		return 0, err
	}
	return int64(binary.Size(v)), nil
}

// WritePaddingTo writes n zero bytes to w.
func WritePaddingTo(w io.Writer, n int) (int64, error) {
	n, err := w.Write(make([]byte, n))
	if err != nil {
		return 0, err
	}
	return int64(n), nil
}
