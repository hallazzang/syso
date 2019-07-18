package versioninfo

import (
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type VersionInfo struct {
	fixedFileInfo  fixedFileInfo
	stringFileInfo *stringFileInfo
	varFileInfo    *varFileInfo
}

func (vi *VersionInfo) SetFileVersion(v uint64) {
	vi.fixedFileInfo.fileVersion = v
}

func (vi *VersionInfo) SetFileVersionString(s string) error {
	v, err := parseVersionString(s)
	if err != nil {
		return errors.Wrap(err, "failed to parse version string")
	}
	vi.fixedFileInfo.fileVersion = v
	return nil
}

func (vi *VersionInfo) SetProductVersion(v uint64) {
	vi.fixedFileInfo.productVersion = v
}

func (vi *VersionInfo) SetProductVersionString(s string) error {
	v, err := parseVersionString(s)
	if err != nil {
		return errors.Wrap(err, "failed to parse version string")
	}
	vi.fixedFileInfo.productVersion = v
	return nil
}

func parseVersionString(s string) (uint64, error) {
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

type fixedFileInfo struct {
	// structVersion uint32 // TODO: do we need it?
	fileVersion    uint64
	productVersion uint64
	fileFlagsMask  uint32
	fileFlags      uint32
	fileOS         uint32
	fileType       uint32
	fileSubtype    uint32
	fileDate       uint64
}

type rawFixedFileInfo struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

type stringFileInfo struct {
	stringTables []stringTable
}

type stringTable struct {
	language uint16
	codepage uint16
	strings  []_string
}

type _string struct {
	key   string
	value string
}

type varFileInfo struct {
	translations []translation
}

type translation struct {
	language uint16
	codepage uint16
}
