package versioninfo

import (
	"encoding/binary"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"unicode/utf16"

	"github.com/pkg/errors"
)

type freezable interface {
	freeze()
	freezableChildren() []freezable
}

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
	Key         [15]uint16
	// Padding     [0]uint16
	Value rawFixedFileInfo
	// Padding2 []uint16
	// Children []interface{}
}

func (vi *VersionInfo) WriteTo(w io.Writer) (int64, error) {
	vi.freeze()
	return 0, nil
}

func (vi *VersionInfo) freeze() {
	for _, c := range vi.freezableChildren() {
		c.freeze()
	}
	vi.valueLength = uint16(binary.Size(rawFixedFileInfo{}))
	vi.length = uint16(binary.Size(rawVersionInfo{}))
	vi.length += paddingLength(vi.length)
	if vi.stringFileInfo != nil {
		vi.length += vi.stringFileInfo.length
	}
	if vi.varFileInfo != nil {
		vi.length += vi.varFileInfo.length
	}
}

func (vi *VersionInfo) freezableChildren() []freezable {
	var r []freezable
	if vi.stringFileInfo != nil {
		r = append(r, vi.stringFileInfo)
	}
	if vi.varFileInfo != nil {
		r = append(r, vi.varFileInfo)
	}
	return r
}

func (vi *VersionInfo) FileVersion() uint64 {
	return vi.fixedFileInfo.fileVersion
}

func (vi *VersionInfo) FileVersionString() string {
	return formatVersionString(vi.fixedFileInfo.fileVersion)
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

func (vi *VersionInfo) ProductVersion() uint64 {
	return vi.fixedFileInfo.productVersion
}

func (vi *VersionInfo) ProductVersionString() string {
	return formatVersionString(vi.fixedFileInfo.productVersion)
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
	Key         [14]uint16
	Padding     [1]uint16
	// Children    []rawStringTable
}

func (sfi *stringFileInfo) freeze() {
	for _, c := range sfi.freezableChildren() {
		c.freeze()
	}
	sfi.length = uint16(binary.Size(rawStringFileInfo{}))
	for _, st := range sfi.stringTables {
		sfi.length += st.length
	}
}

func (sfi *stringFileInfo) freezableChildren() []freezable {
	var r []freezable
	for _, st := range sfi.stringTables {
		r = append(r, st)
	}
	return r
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
	Key         [8]uint16
	Padding     [1]uint16
	// Children    []rawString
}

func (st *stringTable) freeze() {
	for _, c := range st.freezableChildren() {
		c.freeze()
	}
	st.length = uint16(binary.Size(rawStringTable{}))
	for _, s := range st.strings {
		st.length += s.length
	}
}

func (st *stringTable) freezableChildren() []freezable {
	var r []freezable
	for _, s := range st.strings {
		r = append(r, s)
	}
	return r
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
	s.valueLength = uint16(binary.Size(utf16.Encode([]rune(s.value))))
	s.length = uint16(binary.Size(rawString{}))
	s.length += uint16(binary.Size(utf16.Encode([]rune(s.key))))
	s.length += paddingLength(s.length)
	s.length += s.valueLength
}

func (s *_string) freezableChildren() []freezable {
	return nil
}

type varFileInfo struct {
	length uint16
	_var   _var
}

type rawVarFileInfo struct {
	Length      uint16
	ValueLength uint16
	Type        uint16
	Key         [11]uint16
	// Padding     [0]uint16
	Children rawVar
}

func (vfi *varFileInfo) freeze() {
	for _, c := range vfi.freezableChildren() {
		c.freeze()
	}
	vfi.length = uint16(binary.Size(rawVarFileInfo{}))
	vfi.length += vfi._var.length
}

func (vfi *varFileInfo) freezableChildren() []freezable {
	return []freezable{&vfi._var}
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
	Key         [11]uint16
	// Padding     [0]uint16
	// Value []uint32
}

func (v *_var) freeze() {
	v.valueLength = uint16(binary.Size(translation{}) * len(v.translations))
	v.length = uint16(binary.Size(rawVar{}))
	v.length += paddingLength(v.length)
	v.length += v.valueLength
}

func (v *_var) freezableChildren() []freezable {
	return nil
}

type translation struct {
	language uint16
	codepage uint16
}

func paddingLength(n uint16) uint16 {
	return (4 - (n % 4)) % 4
}
