package common

import (
	"encoding/binary"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
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
	if n < 0 {
		return 0, errors.New("n cannot be negative number")
	} else if n == 0 {
		return 0, nil
	}
	n, err := w.Write(make([]byte, n))
	if err != nil {
		return 0, err
	}
	return int64(n), nil
}

// FormatVersionString formats version number in form of "Major.Minor.Patch.Build".
func FormatVersionString(v uint64) string {
	return fmt.Sprintf("%d.%d.%d.%d", (v>>48)&0xffff, (v>>32)&0xffff, (v>>16)&0xffff, v&0xffff)
}

// ParseVersionString parses version string in form of "Major.Minor.Patch.Build" to number.
func ParseVersionString(s string) (uint64, error) {
	r := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)\.(\d+)$`).FindStringSubmatch(s)
	if len(r) == 0 {
		return 0, errors.Errorf("invalid version string format; %q", s)
	}
	var v uint64
	for _, c := range r[1:] {
		n, err := strconv.ParseUint(c, 10, 16)
		if err != nil {
			return 0, errors.Wrapf(err, "failed to parse version component; %q", c)
		}
		v = (v << 16) | n
	}
	return v, nil
}
