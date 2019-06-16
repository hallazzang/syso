package ico

import (
	"encoding/binary"
	"errors"
	"io"
)

// ICO represents an icon group.
type ICO struct {
	dir     *Directory
	entries []*DirectoryEntry
}

// DecodeAll reads an ICO file from r and returns representation
// of the file.
func DecodeAll(r io.Reader) (*ICO, error) {
	var d Directory
	if err := binary.Read(r, binary.LittleEndian, &d); err != nil {
		return nil, err
	}
	if d.Reserved != 0 || d.Type != 1 {
		return nil, errors.New("bad ICO file")
	}

	var entries []*DirectoryEntry
	for i := uint16(0); i < d.Count; i++ {
		var e DirectoryEntry
		if err := binary.Read(r, binary.LittleEndian, &e); err != nil {
			return nil, err
		}
		entries = append(entries, &e)
	}

	return &ICO{
		dir:     &d,
		entries: entries,
	}, nil
}
