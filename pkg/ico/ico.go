package ico

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"

	"github.com/hallazzang/syso/pkg/common"
	"github.com/pkg/errors"
)

// Image represents a single icon image.
type Image struct {
	ID     int // ID must be set manually
	data   []byte
	offset int64 // last read offset; TODO: there should be better way of implementing Read
}

// Read copies image data to p.
func (i *Image) Read(p []byte) (int, error) {
	n := copy(p[:], i.data[i.offset:])
	i.offset += int64(n)
	return n, nil
}

// Size returns image's size.
func (i *Image) Size() int64 {
	return int64(len(i.data))
}

// Group represents an icon group.
type Group struct {
	dir     *directory
	entries []*directoryEntry
	Images  []*Image
}

func (g *Group) Read(p []byte) (int, error) {
	// TODO: is this implementation of Read okay?
	buf := bytes.NewBuffer(p[:0])
	written := 0
	n, err := common.BinaryWriteTo(buf, &groupDirectory{
		Type:  1,
		Count: uint16(len(g.entries)),
	})
	if err != nil {
		return written, errors.Wrap(err, "failed to write icon group directory")
	}
	written += int(n)
	for i, e := range g.entries {
		if g.Images[i].ID == 0 {
			return written, errors.Errorf("image #%d doesn't have an ID", i)
		}
		n, err := common.BinaryWriteTo(buf, &groupDirectoryEntry{
			Width:      e.Width,
			Height:     e.Height,
			ColorCount: e.ColorCount,
			Reserved:   e.Reserved,
			Planes:     e.Planes,
			BitCount:   e.BitCount,
			BytesInRes: e.BytesInRes,
			ID:         uint16(g.Images[i].ID),
		})
		if err != nil {
			return written, errors.Wrapf(err, "failed to write icon group directory entry #%d", i)
		}
		written += int(n)
	}
	return written, nil
}

// Size returns total byte size when g treated as a Blob.
func (g *Group) Size() int64 {
	return int64(binary.Size(&groupDirectory{}) + len(g.entries)*binary.Size(&groupDirectoryEntry{}))
}

// DecodeAll reads an ICO file from r and returns representation
// of the icon group.
func DecodeAll(r Reader) (*Group, error) {
	var d directory
	if err := binary.Read(r, binary.LittleEndian, &d); err != nil {
		return nil, errors.Wrap(err, "failed to read icon directory")
	}
	if d.Reserved != 0 || d.Type != 1 {
		return nil, errors.New("bad ICO file")
	}

	var entries []*directoryEntry
	var images []*Image
	for i := uint16(0); i < d.Count; i++ {
		var e directoryEntry
		if err := binary.Read(r, binary.LittleEndian, &e); err != nil {
			return nil, errors.Wrapf(err, "failed to read icon directory entry #%d", i)
		}
		entries = append(entries, &e)
		data, err := ioutil.ReadAll(io.NewSectionReader(r, int64(e.ImageOffset), int64(e.BytesInRes)))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read icon image #%d's data", i)
		}
		images = append(images, &Image{
			data: data,
		})
	}

	return &Group{
		dir:     &d,
		entries: entries,
		Images:  images,
	}, nil
}
