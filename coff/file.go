// Package coff provides COFF(Common Object File Format)-related
// functionalities.
package coff

import (
	"errors"
	"io"

	"github.com/hallazzang/syso/pkg/common"
)

type rawCOFFFileHeader struct {
	Machine              uint16
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}

type rawCOFFSectionHeader struct {
	Name                 [8]byte
	VirtualSize          uint32
	VirtualAddress       uint32
	SizeOfRawData        uint32
	PointerToRawData     uint32
	PointerToRelocations uint32
	PointerToLineNumbers uint32
	NumberOfRelocations  uint16
	NumberOfLineNumbers  uint16
	Characteristics      uint32
}

type rawCOFFRelocation struct {
	VirtualAddress   uint32
	SymbolTableIndex uint32
	Type             uint16
}

type coffRelocation struct {
	offset uint32
	rawCOFFRelocation
}

type coffSection struct {
	offset uint32
	Section
}

// Section represents a COFF section.
type Section interface {
	Name() string
	WriteTo(w io.Writer) (int64, error)
	Size() int
	Relocations() []Relocation
}

type Relocation interface {
	VirtualAddress() uint32
}

// File is a COFF object file.
type File struct {
	sections    map[string]*coffSection
	symbols     map[string]*coffSymbol
	stringTable map[string]*coffString
}

// New returns newly created file.
func New() *File {
	return &File{}
}

func (f *File) AddSection(section Section) error {
	for _, s := range f.sections {
		if s.Name() == section.Name() {
			return errors.New("duplicate section name")
		}
	}
	// f.sections = append(f.sections, &coffSection{
	// 	Section: section,
	// })
	return nil
}

func (f *File) freeze() {

}

func (f *File) symbolsOffset() uint32 {
	// if len(f.symbols) > 0 {
	// 	return f.symbols[0].offset
	// }
	return 0
}

// WriteTo writes COFF file data to w.
func (f *File) WriteTo(w io.Writer) (int64, error) {
	var written int64

	n, err := common.BinaryWriteTo(w, &rawCOFFFileHeader{
		Machine:              i386Machine,
		NumberOfSections:     uint16(len(f.sections)),
		PointerToSymbolTable: f.symbolsOffset(),
		NumberOfSymbols:      uint32(len(f.symbols)),
		Characteristics:      0x100, // IMAGE_FILE_32BIT_MACHINE
	})
	if err != nil {
		return written, err
	}
	written += n

	for _, s := range f.sections {
		var name [8]byte
		if len(s.Name()) > 8 {
			// TODO: implement
		} else {
			copy(name[:], s.Name())
		}
		n, err := common.BinaryWriteTo(w, &rawCOFFSectionHeader{
			Name:                 name,
			SizeOfRawData:        uint32(s.Size()),
			PointerToRawData:     s.offset,
			PointerToRelocations: 0, // TODO: implement
			NumberOfRelocations:  uint16(len(s.Relocations())),
			Characteristics:      0x40000040, // IMAGE_SCN_MEM_READ|IMAGE_SCN_CNT_INITIALIZED_DATA
		})
		if err != nil {
			return 0, err
		}
		written += n
	}

	return 0, nil
}
