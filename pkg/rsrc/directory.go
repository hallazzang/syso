package rsrc

import (
	"sort"

	"github.com/hallazzang/syso/pkg/common"
	"github.com/pkg/errors"
)

type rawResourceDirectory struct {
	Characteristics     uint32
	TimeDateStamp       uint32
	MajorVersion        uint16
	MinorVersion        uint16
	NumberOfNameEntries uint16
	NumberOfIDEntries   uint16
}

type resourceDirectory struct {
	offset          uint32
	characteristics uint32
	nameEntries     []*resourceDirectoryEntry
	idEntries       []*resourceDirectoryEntry
	strings         map[string]*resourceString
}

func (d *resourceDirectory) addString(s string) *resourceString {
	str, ok := d.strings[s]
	if ok {
		return str
	}
	str = &resourceString{
		string: s,
	}
	d.strings[s] = str
	return str
}

func (d *resourceDirectory) addData(name *string, id *int, blob common.Blob) (*resourceDataEntry, error) {
	e, err := d.addDirectoryEntry(name, id, nil, blob)
	if err != nil {
		return nil, err
	}
	return e.dataEntry, nil
}

func (d *resourceDirectory) addSubdirectory(name *string, id *int, characteristics uint32) (*resourceDirectory, error) {
	e, err := d.addDirectoryEntry(name, id, &characteristics, nil)
	if err != nil {
		return nil, err
	}
	return e.subdirectory, nil
}

func (d *resourceDirectory) addDirectoryEntry(name *string, id *int, characteristics *uint32, blob common.Blob) (*resourceDirectoryEntry, error) {
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
	e := &resourceDirectoryEntry{}
	if name != nil {
		e.name = d.addString(*name)
		d.nameEntries = append(d.nameEntries, e)
	} else {
		e.id = id
		d.idEntries = append(d.idEntries, e)
	}
	if characteristics != nil {
		e.subdirectory = &resourceDirectory{
			characteristics: *characteristics,
			strings:         make(map[string]*resourceString),
		}
	} else {
		e.dataEntry = &resourceDataEntry{
			data: &resourceData{
				Blob: blob,
			},
		}
	}
	d.sort()
	return e, nil
}

func (d *resourceDirectory) sort() {
	sort.SliceStable(d.nameEntries, func(i, j int) bool {
		return d.nameEntries[i].name.string < d.nameEntries[j].name.string
	})
	sort.SliceStable(d.idEntries, func(i, j int) bool {
		return *d.idEntries[i].id < *d.idEntries[j].id
	})
}

func (d *resourceDirectory) walk(cb func(*resourceDirectory) error) error {
	var _walk func(*resourceDirectory) error
	_walk = func(dir *resourceDirectory) error {
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

func (d *resourceDirectory) entries() []*resourceDirectoryEntry {
	return append(append([]*resourceDirectoryEntry{}, d.nameEntries...), d.idEntries...)
}

func (d *resourceDirectory) dataEntries() []*resourceDataEntry {
	var r []*resourceDataEntry
	for _, e := range d.entries() {
		if e.dataEntry != nil {
			r = append(r, e.dataEntry)
		}
	}
	return r
}

func (d *resourceDirectory) subdirectories() []*resourceDirectory {
	var r []*resourceDirectory
	for _, e := range d.entries() {
		if e.subdirectory != nil {
			r = append(r, e.subdirectory)
		}
	}
	return r
}

func (d *resourceDirectory) datas() []*resourceData {
	var r []*resourceData
	for _, e := range d.dataEntries() {
		r = append(r, e.data)
	}
	return r
}

type rawResourceDirectoryEntry struct {
	NameOffsetOrIntegerID               uint32
	DataEntryOffsetOrSubdirectoryOffset uint32
}

type resourceDirectoryEntry struct {
	offset       uint32
	name         *resourceString
	id           *int
	dataEntry    *resourceDataEntry
	subdirectory *resourceDirectory
}

type resourceString struct {
	offset uint32
	string
}

type rawResourceDataEntry struct {
	DataRVA  uint32
	Size     uint32
	Codepage uint32
	Reserved uint32
}

type resourceDataEntry struct {
	offset uint32
	data   *resourceData
}

type resourceData struct {
	offset uint32
	common.Blob
}
