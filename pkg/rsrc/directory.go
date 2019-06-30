package rsrc

import (
	"sort"

	"github.com/hallazzang/syso/pkg/common"
	"github.com/pkg/errors"
)

type rawDirectory struct {
	Characteristics     uint32
	TimeDateStamp       uint32
	MajorVersion        uint16
	MinorVersion        uint16
	NumberOfNameEntries uint16
	NumberOfIDEntries   uint16
}

type Directory struct {
	offset          uint32
	characteristics uint32
	nameEntries     []*DirectoryEntry
	idEntries       []*DirectoryEntry
	strings         map[string]*String
}

func (d *Directory) addString(s string) *String {
	str, ok := d.strings[s]
	if ok {
		return str
	}
	str = &String{
		string: s,
	}
	d.strings[s] = str
	return str
}

func (d *Directory) addData(name *string, id *int, blob common.Blob) (*DataEntry, error) {
	e, err := d.addDirectoryEntry(name, id, nil, blob)
	if err != nil {
		return nil, err
	}
	return e.dataEntry, nil
}

func (d *Directory) addSubdirectory(name *string, id *int, characteristics uint32) (*Directory, error) {
	e, err := d.addDirectoryEntry(name, id, &characteristics, nil)
	if err != nil {
		return nil, err
	}
	return e.subdirectory, nil
}

func (d *Directory) addDirectoryEntry(name *string, id *int, characteristics *uint32, blob common.Blob) (*DirectoryEntry, error) {
	for _, e := range d.entries() {
		if name != nil {
			if e.name != nil && e.name.string == *name {
				return nil, errors.New("duplicate directory entry name")
			}
		} else {
			if e.id != nil && *e.id == *id {
				return nil, errors.New("duplicate directory entry id")
			}
		}
	}
	e := &DirectoryEntry{}
	if name != nil {
		e.name = d.addString(*name)
		d.nameEntries = append(d.nameEntries, e)
	} else {
		e.id = id
		d.idEntries = append(d.idEntries, e)
	}
	if characteristics != nil {
		e.subdirectory = &Directory{
			characteristics: *characteristics,
			strings:         make(map[string]*String),
		}
	} else {
		e.dataEntry = &DataEntry{
			data: &Data{
				Blob: blob,
			},
		}
	}
	d.sort()
	return e, nil
}

func (d *Directory) sort() {
	sort.SliceStable(d.nameEntries, func(i, j int) bool {
		return d.nameEntries[i].name.string < d.nameEntries[j].name.string
	})
	sort.SliceStable(d.idEntries, func(i, j int) bool {
		return *d.idEntries[i].id < *d.idEntries[j].id
	})
}

func (d *Directory) walk(cb func(*Directory) error) error {
	var _walk func(*Directory) error
	_walk = func(dir *Directory) error {
		if err := cb(dir); err != nil {
			return err
		}
		for _, subdir := range dir.subdirectories() {
			if err := _walk(subdir); err != nil {
				return err
			}
		}
		return nil
	}

	return _walk(d)
}

func (d *Directory) entries() []*DirectoryEntry {
	return append(append([]*DirectoryEntry{}, d.nameEntries...), d.idEntries...)
}

func (d *Directory) dataEntries() []*DataEntry {
	var r []*DataEntry
	for _, e := range d.entries() {
		if e.dataEntry != nil {
			r = append(r, e.dataEntry)
		}
	}
	return r
}

func (d *Directory) subdirectories() []*Directory {
	var r []*Directory
	for _, e := range d.entries() {
		if e.subdirectory != nil {
			r = append(r, e.subdirectory)
		}
	}
	return r
}

func (d *Directory) datas() []*Data {
	var r []*Data
	for _, e := range d.dataEntries() {
		r = append(r, e.data)
	}
	return r
}

type rawDirectoryEntry struct {
	NameOffsetOrIntegerID               uint32
	DataEntryOffsetOrSubdirectoryOffset uint32
}

type DirectoryEntry struct {
	offset       uint32
	name         *String
	id           *int
	dataEntry    *DataEntry
	subdirectory *Directory
}

type String struct {
	offset uint32
	string
}
