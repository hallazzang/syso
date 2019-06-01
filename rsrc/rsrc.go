package rsrc

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf16"

	"github.com/hallazzang/syso/internal/common"
)

// Section is a .rsrc section.
type Section struct {
	rootDir *directory
}

// New returns newly created .rsrc section.
func New() *Section {
	return &Section{
		rootDir: &directory{},
	}
}

// Name returns section's name, which is .rsrc in general.
func (s *Section) Name() string {
	return ".rsrc"
}

func (s *Section) freeze() error {
	var offset uint32
	_ = offset

	// TODO: work on this function. added offset fields to each structures

	if err := s.rootDir.walk(func(dir *directory) error {
		offset += uint32(binary.Size(&resourceDirectory{}))
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// WriteTo writes section data to w.
func (s *Section) WriteTo(w io.Writer) (int64, error) {
	var written int64

	if err := s.rootDir.walk(func(dir *directory) error {
		fmt.Printf("%08X %d\n", written, dir.characteristics)
		n, err := common.BinaryWriteTo(w, &resourceDirectory{
			Characteristics:     dir.characteristics,
			NumberOfNameEntries: uint16(len(dir.nameEntries)),
			NumberOfIDEntries:   uint16(len(dir.idEntries)),
		})
		if err != nil {
			return err
		}
		written += n

		for _, e := range dir.nameEntries {
			// TODO: implement
			_ = e
		}
		for _, e := range dir.idEntries {
			// TODO: implement
			_ = e
		}

		return nil
	}); err != nil {
		return written, err
	}

	if err := s.rootDir.walk(func(dir *directory) error {
		for _, e := range dir.nameEntries {
			s := utf16.Encode([]rune(*e.name))
			n, err := common.BinaryWriteTo(w, uint16(len(s)))
			if err != nil {
				return err
			}
			written += n
			n, err = common.BinaryWriteTo(w, s)
			if err != nil {
				return err
			}
			written += n
		}
		return nil
	}); err != nil {
		return written, err
	}

	if err := s.rootDir.walk(func(dir *directory) error {
		return nil
	}); err != nil {
		return written, err
	}

	return written, nil
}

// AddIconByID adds an icon resource identified by an integer id
// to section.
func (s *Section) AddIconByID(id int, data common.Blob) error {
	panic("not implemented")
	return nil
}

// AddIconByName adds an icon resource identified by name to section.
func (s *Section) AddIconByName(name string, data common.Blob) error {
	panic("not implemented")
	return nil
}
