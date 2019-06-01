// Package coff provides COFF(Common Object File Format)-related
// functionalities.
package coff

// Section represents a COFF section.
type Section interface {
	Name() string
}

type coffFileHeader struct {
	Machine              uint16
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}

type coffSectionHeader struct {
	Name                 [8]uint8
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

// File is a COFF object file.
type File struct {
}

// New returns newly created COFF object file.
func New() *File {
	panic("not implemented")
	return &File{}
}
