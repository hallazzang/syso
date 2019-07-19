package versioninfo

import "testing"

func TestFormatVersionString(t *testing.T) {
	if formatVersionString(0x0001000200030004) != "1.2.3.4" {
		t.Errorf("failed")
	}
	vi := &VersionInfo{}
	if err := vi.SetFileVersionString("1.2.3.4"); err != nil {
		t.Fatal(err)
	}
	if vi.FileVersionString() != "1.2.3.4" {
		t.Errorf("failed")
	}
}

func TestParseVersionString(t *testing.T) {
	v, err := parseVersionString("1.2.3.4")
	if err != nil {
		t.Fatal(err)
	}
	if v != 0x0001000200030004 {
		t.Errorf("mismatching version; expected 0x0001000200030004, got %#016x", v)
	}
}
