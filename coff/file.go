// Package coff provides COFF(Common Object File Format)-related
// functionalities.
package coff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/hallazzang/syso/pkg/common"
)

type rawFileHeader struct {
	Machine              uint16
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}

type File struct {
	sections        []*section
	symbolsOffset   uint32
	strings         []*_string
	stringTable     map[string]*_string
	stringTableSize uint32
}

func New() *File {
	return &File{
		stringTable: make(map[string]*_string),
	}
}

func (f *File) AddSection(s Section) error {
	for _, sec := range f.sections {
		if sec.Name() == s.Name() {
			return errors.New("duplicate section name")
		}
	}
	f.sections = append(f.sections, &section{
		Section: s,
	})
	if len(s.Name()) > 8 {
		if _, ok := f.stringTable[s.Name()]; !ok {
			str := &_string{
				b: []byte(s.Name() + "\x00"), // null-terminated UTF-8 encoded string
			}
			f.strings = append(f.strings, str)
			f.stringTable[s.Name()] = str
		}
	}
	return nil
}

func (f *File) freeze() {
	offset := uint32(binary.Size(&rawFileHeader{}))
	offset += uint32(binary.Size(&rawSectionHeader{}) * len(f.sections))
	for _, s := range f.sections {
		s.dataOffset = offset
		offset += uint32(s.Size())
	}
	for _, s := range f.sections {
		s.relocationsOffset = offset
		offset += uint32(binary.Size(&rawRelocation{}) * len(s.Relocations()))
	}
	f.symbolsOffset = offset
	offset += uint32(binary.Size(&rawSymbol{}) * len(f.sections))
	offset += 4  // string table size
	so := offset // start offset of string table
	for _, s := range f.strings {
		s.offset = offset
		offset += uint32(len(s.b))
	}
	f.stringTableSize = offset - so + 4
}

// WriteTo writes COFF file data to w.
func (f *File) WriteTo(w io.Writer) (int64, error) {
	var written int64

	f.freeze()

	n, err := common.BinaryWriteTo(w, &rawFileHeader{
		Machine:              0x14c, // IMAGE_FILE_MACHINE_I386
		NumberOfSections:     uint16(len(f.sections)),
		PointerToSymbolTable: f.symbolsOffset,
		NumberOfSymbols:      uint32(len(f.sections)),
		Characteristics:      0x0100, // IMAGE_FILE_32BIT_MACHINE
	})
	if err != nil {
		return written, err
	}
	written += n

	for _, s := range f.sections {
		var name [8]byte
		if len(s.Name()) > 8 {
			copy(name[:], fmt.Sprintf("/%d", f.stringTable[s.Name()].offset))
		} else {
			copy(name[:], s.Name())
		}
		n, err := common.BinaryWriteTo(w, &rawSectionHeader{
			Name:                 name,
			SizeOfRawData:        uint32(s.Size()),
			PointerToRawData:     s.dataOffset,
			PointerToRelocations: s.relocationsOffset,
			NumberOfRelocations:  uint16(len(s.Relocations())),
			Characteristics:      0x40000040, // IMAGE_SCN_MEM_READ|IMAGE_SCN_CNT_INITIALIZED_DATA
		})
		if err != nil {
			return written, err
		}
		written += n
	}

	for _, s := range f.sections {
		n, err := s.WriteTo(w)
		if err != nil {
			return written, err
		}
		written += n
	}

	for i, s := range f.sections {
		for _, r := range s.Relocations() {
			n, err := common.BinaryWriteTo(w, &rawRelocation{
				VirtualAddress:   r.VirtualAddress(),
				SymbolTableIndex: uint32(i),
				Type:             0x0007, // IMAGE_REL_I386_DIR32NB
			})
			if err != nil {
				return written, err
			}
			written += n
		}
	}

	for i, s := range f.sections {
		var name [8]byte
		if len(s.Name()) > 8 {
			binary.LittleEndian.PutUint32(name[4:], f.stringTable[s.Name()].offset)
		} else {
			copy(name[:], s.Name())
		}
		n, err := common.BinaryWriteTo(w, &rawSymbol{
			Name:          name,
			SectionNumber: uint16(i) + 1,
			StorageClass:  3, // IMAGE_SYM_CLASS_STATIC
		})
		if err != nil {
			return written, err
		}
		written += n
	}

	n, err = common.BinaryWriteTo(w, f.stringTableSize)
	if err != nil {
		return written, err
	}
	written += n
	for _, s := range f.strings {
		n, err := common.BinaryWriteTo(w, s.b)
		if err != nil {
			return written, err
		}
		written += n
	}

	return written, nil
}
