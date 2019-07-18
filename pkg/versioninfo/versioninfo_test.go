package versioninfo

import "testing"

func TestParseVersionString(t *testing.T) {
	v, err := parseVersionString("1.2.3.4")
	if err != nil {
		t.Fatal(err)
	}
	if v != 0x0001000200030004 {
		t.Errorf("mismatching version; expected 0x0001000200030004, got %#016x", v)
	}
}
