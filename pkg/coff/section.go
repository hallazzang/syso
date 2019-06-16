package coff

import "io"

type rawSectionHeader struct {
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

type section struct {
	dataOffset        uint32
	relocationsOffset uint32
	Section
}

// Section is a COFF section.
type Section interface {
	Name() string
	WriteTo(w io.Writer) (int64, error)
	Size() int
	Relocations() []Relocation
}
