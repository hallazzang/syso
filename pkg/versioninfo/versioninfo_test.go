package versioninfo

import (
	"bytes"
	"testing"
)

func TestFormatVersionString(t *testing.T) {
	vi := &VersionInfo{}
	if err := vi.SetFileVersionString("1.2.3.4"); err != nil {
		t.Fatal(err)
	}
	if vi.FileVersionString() != "1.2.3.4" {
		t.Errorf("failed")
	}
}

func TestString(t *testing.T) {
	vi := &VersionInfo{}
	vi.SetString(0x0409, 0x04b0, "foo", "bar")
	if s, ok := vi.String(0x0409, 0x04b0, "foo"); !ok {
		t.Fatal("cannot get string")
	} else if s != "bar" {
		t.Fatal("wrong string")
	}
	if _, ok := vi.String(0x1, 0x2, "foo"); ok {
		t.Fatal("must not get string")
	}
}

func TestFreezeEmpty(t *testing.T) {
	vi := &VersionInfo{}
	vi.freeze()
	if vi.length != 92 {
		t.Errorf("wrong VersionInfo.length; expected 88, got %d", vi.length)
	}
	if vi.valueLength != 52 {
		t.Errorf("wrong VersionInfo.valueLength; expected 52, got %d", vi.valueLength)
	}
}

func TestFreeze(t *testing.T) {
	vi := &VersionInfo{}
	vi.SetString(0x0409, 0x04b0, "foo", "bar")
	vi.freeze()
	if vi.length != 176 {
		t.Errorf("wrong VersionInfo.length; expected 176, got %d", vi.length)
	}
	if vi.stringFileInfo.length != 84 {
		t.Errorf("wrong VersionInfo.stringFileInfo.length; expected 84, got %d", vi.stringFileInfo.length)
	}
	if vi.stringFileInfo.stringTables[0].length != 48 {
		t.Errorf("wrong VersionInfo.stringFileInfo.stringTables[0].length; expected 48, got %d", vi.stringFileInfo.stringTables[0].length)
	}
	if vi.stringFileInfo.stringTables[0].strings[0].length != 24 {
		t.Errorf("wrong VersionInfo.stringFileInfo.stringTables[0].strings[0].length; expected 24, got %d", vi.stringFileInfo.stringTables[0].strings[0].length)
	}
	if vi.stringFileInfo.stringTables[0].strings[0].valueLength != 4 {
		t.Errorf("wrong VersionInfo.stringFileInfo.stringTables[0].strings[0].valueLength; expected 4, got %d", vi.stringFileInfo.stringTables[0].strings[0].valueLength)
	}
}

func TestWrite(t *testing.T) {
	vi := &VersionInfo{}
	vi.SetString(0x0409, 0x04b0, "foo", "bar")

	b := new(bytes.Buffer)
	n, err := vi.WriteTo(b)
	if err != nil {
		t.Fatal(err)
	}
	if n != int64(b.Len()) || n != int64(vi.length) {
		t.Fatal("wrong length")
	}
}
