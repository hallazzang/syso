package common

import (
	"bytes"
	"io"
	"testing"

	"github.com/pkg/errors"
)

type badWriter struct{}

func (*badWriter) Write([]byte) (int, error) {
	return 0, errors.New("you cannot write to bad writer")
}

func TestBinaryWriteTo(t *testing.T) {
	for _, tc := range []struct {
		Value  interface{}
		Length int
		Result []byte
	}{
		{byte(0), 1, []byte{0}},
		{int32(0x12345678), 4, []byte{0x78, 0x56, 0x34, 0x12}},
	} {
		b := &bytes.Buffer{}
		n, err := BinaryWriteTo(b, tc.Value)
		if err != nil {
			t.Fatal(err)
		}
		if int(n) != tc.Length {
			t.Fatalf("wrong written length; expected %d, got %d", tc.Length, n)
		}
		if r := b.Bytes(); !bytes.Equal(r, tc.Result) {
			t.Fatalf("wrong result; expected %+q, got %+q", tc.Result, r)
		}
	}
}

func TestBinaryWriteTo_invalidValue(t *testing.T) {
	for _, tc := range []interface{}{
		0,
		"hello",
		struct{ Value int }{1},
	} {
		b := &bytes.Buffer{}
		n, err := BinaryWriteTo(b, tc)
		if err == nil {
			t.Fatalf("expected failure for value %v(type %T), got no error", tc, tc)
		}
		if n != 0 {
			t.Fatalf("wrong written length; expected 0, got %d", n)
		}
	}
}

func TestBinaryWriteTo_invalidWriter(t *testing.T) {
	for _, tc := range []io.Writer{
		&badWriter{},
	} {
		if _, err := BinaryWriteTo(tc, byte(0)); err == nil {
			t.Fatalf("expected failure for writer type %T, got no error", tc)
		}
	}
}

func TestWritePaddingTo(t *testing.T) {
	for _, tc := range []int{
		0,
		1,
		100,
	} {
		b := &bytes.Buffer{}
		n, err := WritePaddingTo(b, tc)
		if err != nil {
			t.Fatal(err)
		}
		if int(n) != tc {
			t.Fatalf("wrong written length; expected %d, got %d", tc, n)
		}
		if b.Len() != tc {
			t.Fatalf("wrong result length; expected %d, got %d", tc, b.Len())
		}
		for _, c := range b.Bytes() {
			if c != 0 {
				t.Fatal("result contains byte other than 0")
			}
		}
	}
}

func TestWritePaddingTo_invalidValue(t *testing.T) {
	for _, tc := range []int{
		-1,
	} {
		b := &bytes.Buffer{}
		n, err := WritePaddingTo(b, tc)
		if err == nil {
			t.Fatal("expected failure, got no error")
		}
		if n != 0 {
			t.Fatalf("wrong written length; expected 0, got %d", n)
		}
	}
}

func TestWritePaddingTo_invalidWriter(t *testing.T) {
	for _, tc := range []io.Writer{
		&badWriter{},
	} {
		if _, err := WritePaddingTo(tc, 1); err == nil {
			t.Fatalf("expected failure for writer type %T, got no error", tc)
		}
	}
}

func TestFormatVersionString(t *testing.T) {
	for _, tc := range []struct {
		Value  uint64
		Result string
	}{
		{0, "0.0.0.0"},
		{1, "0.0.0.1"},
		{256, "0.0.0.256"},
		{65535, "0.0.0.65535"},
		{65536, "0.0.1.0"},
		{0xffffffff00000000, "65535.65535.0.0"},
		{0x0001000200030004, "1.2.3.4"},
	} {
		if r := FormatVersionString(tc.Value); r != tc.Result {
			t.Fatalf("wrong result; expected %+q, got %+q", tc.Result, r)
		}
	}
}

func TestParseVersionString(t *testing.T) {
	for _, tc := range []struct {
		Value  string
		Result uint64
	}{
		{"0.0.0.0", 0},
		{"0.0.0.1", 1},
		{"0.0.0.256", 256},
		{"0.0.0.65535", 65535},
		{"0.0.1.0", 65536},
		{"65535.65535.0.0", 0xffffffff00000000},
		{"1.2.3.4", 0x0001000200030004},
	} {
		r, err := ParseVersionString(tc.Value)
		if err != nil {
			t.Fatal(err)
		}
		if r != tc.Result {
			t.Fatalf("wrong result; expected %d, got %d", tc.Result, r)
		}
	}
}

func TestParseVersionString_invalidValue(t *testing.T) {
	for _, tc := range []string{
		"",
		"1",
		"1.2.3.",
		"-1.-1.-1.-1",
		"65536.65536.65536.65536",
		"a.b.c.d",
	} {
		r, err := ParseVersionString(tc)
		if err == nil {
			t.Fatalf("expected failure for %q, got no error", tc)
		}
		if r != 0 {
			t.Fatalf("wrong result; expected 0, got %d", r)
		}
	}
}
