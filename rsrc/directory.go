package rsrc

import (
	"sort"
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
	strings         []*resourceString
}

func (d *resourceDirectory) addDataEntryByName(name string, blob Blob) {
	// TODO: check for duplicate name
	nameString := &resourceString{
		string: name,
	}
	d.strings = append(d.strings, nameString)
	d.nameEntries = append(d.nameEntries, &resourceDirectoryEntry{
		name: nameString,
		dataEntry: &resourceDataEntry{
			data: &resourceData{
				blob: blob,
			},
		},
	})
	d.sort()
}

func (d *resourceDirectory) addDataEntryByID(id int, blob Blob) {
	d.idEntries = append(d.idEntries, &resourceDirectoryEntry{
		id: &id,
		dataEntry: &resourceDataEntry{
			data: &resourceData{
				blob: blob,
			},
		},
	})
	d.sort()
}

func (d *resourceDirectory) addSubdirectoryByName(name string, characteristics uint32) *resourceDirectory {
	// TODO: check for duplicate name
	nameString := &resourceString{
		string: name,
	}
	d.strings = append(d.strings, nameString)
	subdir := &resourceDirectory{
		characteristics: characteristics,
	}
	d.nameEntries = append(d.nameEntries, &resourceDirectoryEntry{
		name:         nameString,
		subdirectory: subdir,
	})
	d.sort()
	return subdir
}

func (d *resourceDirectory) addSubdirectoryByID(id int, characteristics uint32) *resourceDirectory {
	subdir := &resourceDirectory{
		characteristics: characteristics,
	}
	d.idEntries = append(d.idEntries, &resourceDirectoryEntry{
		id:           &id,
		subdirectory: subdir,
	})
	d.sort()
	return subdir
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

	if err := _walk(d); err != nil {
		return err
	}
	return nil
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
	blob   Blob
}
