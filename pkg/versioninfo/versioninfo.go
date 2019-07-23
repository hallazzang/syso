package versioninfo

import (
	"encoding/binary"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"unicode/utf16"

	"github.com/hallazzang/syso/pkg/common"
	"github.com/pkg/errors"
)

// VersionInfo is a root structure for a version info resource.
type VersionInfo struct {
	length         uint16
	valueLength    uint16
	fixedFileInfo  fixedFileInfo // TODO: make fixedFileInfo optional
	stringFileInfo *stringFileInfo
	varFileInfo    *varFileInfo
}

type rawVersionInfo struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	Key         [16]uint16
	Padding     [2]byte
	Value       rawFixedFileInfo
	// Padding2    []uint16
	// Children    []interface{}
}

// New returns new version info resource with some fields set to default.
func New() *VersionInfo {
	return &VersionInfo{
		fixedFileInfo: fixedFileInfo{
			fileFlagsMask: 0x0000003f, // TODO: find out what it means
			fileOS:        0x00040004,
			fileType:      0x00000001,
		},
	}
}

func (vi *VersionInfo) freeze() {
	if vi.stringFileInfo != nil {
		vi.stringFileInfo.freeze()
	}
	if vi.varFileInfo != nil {
		vi.varFileInfo.freeze()
	}

	vi.valueLength = uint16(binary.Size(rawFixedFileInfo{}))
	vi.length = uint16(binary.Size(rawVersionInfo{}))
	if vi.stringFileInfo != nil {
		vi.length += vi.stringFileInfo.length
	}
	vi.length += paddingLength(vi.length)
	if vi.varFileInfo != nil {
		vi.length += vi.varFileInfo.length
	}
}

// WriteTo writes content of VersionInfo to w in binary format.
func (vi *VersionInfo) WriteTo(w io.Writer) (int64, error) {
	vi.freeze()

	written, err := common.BinaryWriteTo(w, rawVersionInfo{
		Length:      vi.length,
		ValueLength: vi.valueLength,
		Type:        0,
		Key:         [16]uint16{0x56, 0x53, 0x5f, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x00}, // L"VS_VERSION_INFO"
		Value: rawFixedFileInfo{
			Signature:        0xFEEF04BD,
			StrucVersion:     0x00010000, // TODO: needed?
			FileVersionMS:    uint32(vi.fixedFileInfo.fileVersion >> 32),
			FileVersionLS:    uint32(vi.fixedFileInfo.fileVersion & 0xffffffff),
			ProductVersionMS: uint32(vi.fixedFileInfo.productVersion >> 32),
			ProductVersionLS: uint32(vi.fixedFileInfo.productVersion & 0xffffffff),
			FileFlagsMask:    vi.fixedFileInfo.fileFlagsMask,
			FileFlags:        vi.fixedFileInfo.fileFlags,
			FileOS:           vi.fixedFileInfo.fileOS,
			FileType:         vi.fixedFileInfo.fileType,
			FileSubtype:      vi.fixedFileInfo.fileSubtype,
			FileDateMS:       uint32(vi.fixedFileInfo.fileDate >> 32),
			FileDateLS:       uint32(vi.fixedFileInfo.fileDate & 0xffffffff),
		},
	})
	if err != nil {
		return written, err
	}

	if vi.stringFileInfo != nil {
		n, err := vi.stringFileInfo.WriteTo(w)
		if err != nil {
			return written, err
		}
		written += n
	}

	if padding := paddingLength(uint16(written)); padding > 0 { // TODO: make paddingLength to accept int64
		n, err := common.WritePaddingTo(w, int(padding))
		if err != nil {
			return written, err
		}
		written += n
	}

	if vi.varFileInfo != nil {
		n, err := vi.varFileInfo.WriteTo(w)
		if err != nil {
			return written, err
		}
		written += n
	}

	return written, nil
}

// FileVersion returns file version in integer.
func (vi *VersionInfo) FileVersion() uint64 {
	return vi.fixedFileInfo.fileVersion
}

// FileVersionString returns file version in string, formatted in
// "Major.Minor.Patch.Build" form.
func (vi *VersionInfo) FileVersionString() string {
	return formatVersionString(vi.fixedFileInfo.fileVersion)
}

// SetFileVersion sets file version in integer.
func (vi *VersionInfo) SetFileVersion(v uint64) {
	vi.fixedFileInfo.fileVersion = v
}

// SetFileVersionString sets file version in string, returns error
// if s is not in a form of "Major.Minor.Patch.Build".
func (vi *VersionInfo) SetFileVersionString(s string) error {
	v, err := parseVersionString(s)
	if err != nil {
		return errors.Wrap(err, "failed to parse version string")
	}
	vi.fixedFileInfo.fileVersion = v
	return nil
}

// ProductVersion returns product version in integer.
func (vi *VersionInfo) ProductVersion() uint64 {
	return vi.fixedFileInfo.productVersion
}

// ProductVersionString returns product version in string, formatted in
// "Major.Minor.Patch.Build" form.
func (vi *VersionInfo) ProductVersionString() string {
	return formatVersionString(vi.fixedFileInfo.productVersion)
}

// SetProductVersion sets product version in integer.
func (vi *VersionInfo) SetProductVersion(v uint64) {
	vi.fixedFileInfo.productVersion = v
}

// SetProductVersionString sets product version in string, returns error
// if s is not in a form "Major.Minor.Patch.Build".
func (vi *VersionInfo) SetProductVersionString(s string) error {
	v, err := parseVersionString(s)
	if err != nil {
		return errors.Wrap(err, "failed to parse version string")
	}
	vi.fixedFileInfo.productVersion = v
	return nil
}

func formatVersionString(v uint64) string {
	return fmt.Sprintf("%d.%d.%d.%d", (v>>48)&0xffff, (v>>32)&0xffff, (v>>16)&0xffff, v&0xffff)
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

// String returns a string value from string table which is indicated
// by given language and codepage pair.
// The second return value indicates whether the key has found or not.
func (vi *VersionInfo) String(language, codepage uint16, key string) (string, bool) {
	st := vi.stringTable(language, codepage, false)
	if st == nil {
		return "", false
	}
	for _, s := range st.strings {
		if s.key == key {
			return s.value, true
		}
	}
	return "", false
}

// SetString sets a string value for a key in string table which is
// indicated by given language and codepage pair.
func (vi *VersionInfo) SetString(language, codepage uint16, key, value string) {
	st := vi.stringTable(language, codepage, true)
	f := false
	for _, s := range st.strings {
		if s.key == key {
			s.value = value
			f = true
			break
		}
	}
	if !f {
		st.strings = append(st.strings, &_string{
			key:   key,
			value: value,
		})
	}
}

func (vi *VersionInfo) stringTable(language, codepage uint16, createIfNotExists bool) *stringTable {
	if vi.stringFileInfo == nil {
		if !createIfNotExists {
			return nil
		}
		vi.stringFileInfo = &stringFileInfo{}
	}
	var st *stringTable
	for _, t := range vi.stringFileInfo.stringTables {
		if t.language == language && t.codepage == codepage {
			st = t
			break
		}
	}
	if st == nil {
		if !createIfNotExists {
			return nil
		}
		st = &stringTable{
			language: language,
			codepage: codepage,
		}
		vi.stringFileInfo.stringTables = append(vi.stringFileInfo.stringTables, st)
	}
	return st
}

// AddTranslation adds a translation info.
func (vi *VersionInfo) AddTranslation(language, codepage uint16) {
	if vi.varFileInfo == nil {
		vi.varFileInfo = &varFileInfo{}
	}
	for _, t := range vi.varFileInfo._var.translations {
		if t.language == language && t.codepage == codepage {
			return
		}
	}
	vi.varFileInfo._var.translations = append(vi.varFileInfo._var.translations, &translation{
		language: language,
		codepage: codepage,
	})
}

// TODO: add methods for getting/setting FileFlags, OS, etc.

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
	length       uint16
	stringTables []*stringTable
}

type rawStringFileInfo struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	Key         [15]uint16
	// Padding     [0]byte
	// Children    []rawStringTable
}

func (sfi *stringFileInfo) freeze() {
	for _, st := range sfi.stringTables {
		st.freeze()
	}

	sfi.length = uint16(binary.Size(rawStringFileInfo{}))
	for _, st := range sfi.stringTables {
		sfi.length += paddingLength(sfi.length)
		sfi.length += st.length
	}
}

func (sfi *stringFileInfo) WriteTo(w io.Writer) (int64, error) {
	written, err := common.BinaryWriteTo(w, rawStringFileInfo{
		Length:      sfi.length,
		ValueLength: 0,
		Type:        1,
		Key:         [15]uint16{0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x00}, // L"StringFileInfo"
	})
	if err != nil {
		return written, err
	}

	for _, st := range sfi.stringTables {
		if padding := paddingLength(uint16(written)); padding > 0 {
			n, err := common.WritePaddingTo(w, int(padding))
			if err != nil {
				return written, err
			}
			written += n
		}

		n, err := st.WriteTo(w)
		if err != nil {
			return written, err
		}
		written += n
	}

	return written, nil
}

type stringTable struct {
	length   uint16
	language uint16
	codepage uint16
	strings  []*_string
}

type rawStringTable struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	Key         [9]uint16
	// Padding     [0]byte
	// Children    []rawString
}

func (st *stringTable) freeze() {
	for _, s := range st.strings {
		s.freeze()
	}

	st.length = uint16(binary.Size(rawStringTable{}))
	for _, s := range st.strings {
		st.length += paddingLength(st.length)
		st.length += s.length
	}
}

func (st *stringTable) WriteTo(w io.Writer) (int64, error) {
	var key [9]uint16
	copy(key[:], utf16.Encode([]rune(fmt.Sprintf("%04X%04X", st.language, st.codepage))))
	written, err := common.BinaryWriteTo(w, rawStringTable{
		Length:      st.length,
		ValueLength: 0,
		Type:        1,
		Key:         key,
	})
	if err != nil {
		return written, err
	}

	for _, s := range st.strings {
		if padding := paddingLength(uint16(written)); padding > 0 {
			n, err := common.WritePaddingTo(w, int(padding))
			if err != nil {
				return written, err
			}
			written += n
		}

		n, err := s.WriteTo(w)
		if err != nil {
			return written, err
		}
		written += n
	}

	return written, nil
}

type _string struct {
	length      uint16
	valueLength uint16
	key         string
	value       string
}

type rawString struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	// Key         []uint16
	// Padding     []uint16
	// Value       []uint16
}

func (s *_string) freeze() {
	s.valueLength = uint16(binary.Size(utf16.Encode([]rune(s.value+"\x00")))) / 2
	s.length = uint16(binary.Size(rawString{}))
	s.length += uint16(binary.Size(utf16.Encode([]rune(s.key + "\x00"))))
	s.length += paddingLength(s.length)
	s.length += s.valueLength * 2
}

func (s *_string) WriteTo(w io.Writer) (int64, error) {
	written, err := common.BinaryWriteTo(w, rawString{
		Length:      s.length,
		ValueLength: s.valueLength,
		Type:        1,
	})
	if err != nil {
		return written, err
	}

	n, err := common.BinaryWriteTo(w, utf16.Encode([]rune(s.key+"\x00")))
	if err != nil {
		return written, err
	}
	written += n

	n, err = common.WritePaddingTo(w, int(paddingLength(uint16(written))))
	if err != nil {
		return written, err
	}
	written += n

	n, err = common.BinaryWriteTo(w, utf16.Encode([]rune(s.value+"\x00")))
	if err != nil {
		return written, err
	}
	written += n

	return written, nil
}

type varFileInfo struct {
	length uint16
	_var   _var
}

type rawVarFileInfo struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	Key         [12]uint16
	Padding     [2]byte
	// Children rawVar
}

func (vfi *varFileInfo) freeze() {
	vfi._var.freeze()

	vfi.length = uint16(binary.Size(rawVarFileInfo{}))
	vfi.length += vfi._var.length
}

func (vfi *varFileInfo) WriteTo(w io.Writer) (int64, error) {
	written, err := common.BinaryWriteTo(w, rawVarFileInfo{
		Length:      vfi.length,
		ValueLength: 0,
		Type:        1,
		Key:         [12]uint16{0x56, 0x61, 0x72, 0x46, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x00}, // L"VarFileInfo"
	})
	if err != nil {
		return written, err
	}

	n, err := vfi._var.WriteTo(w)
	if err != nil {
		return written, err
	}
	written += n

	return written, nil
}

type _var struct {
	length       uint16
	valueLength  uint16
	translations []*translation
}

type rawVar struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	Key         [12]uint16
	Padding     [2]byte
	// Value       []uint32
}

func (v *_var) freeze() {
	v.valueLength = uint16(binary.Size(translation{}) * len(v.translations))
	v.length = uint16(binary.Size(rawVar{}))
	v.length += v.valueLength
}

func (v *_var) WriteTo(w io.Writer) (int64, error) {
	written, err := common.BinaryWriteTo(w, rawVar{
		Length:      v.length,
		ValueLength: v.valueLength,
		Type:        0,
		Key:         [12]uint16{0x54, 0x72, 0x61, 0x6e, 0x73, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x00}, // L"Translation"
	})
	if err != nil {
		return written, err
	}

	for _, t := range v.translations {
		n, err := common.BinaryWriteTo(w, t.language)
		if err != nil {
			return written, err
		}
		written += n

		n, err = common.BinaryWriteTo(w, t.codepage)
		if err != nil {
			return written, err
		}
		written += n
	}

	return written, nil
}

type translation struct {
	language uint16
	codepage uint16
}

// paddingLength returns how many bytes to pad when there are n bytes already
// to fit the data on 32-bit boundary.
func paddingLength(n uint16) uint16 {
	// kinda weird calculation, but it works
	// p := uint16(1)
	// for (n+p)%4 != 0 {
	// 	p++
	// }
	// return p
	return (4 - (n % 4)) % 4 // this was my first attempt
}
