package rsrc

import (
	"sort"

	"github.com/hallazzang/syso/internal/common"
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
	strings         []string
}

func (d *resourceDirectory) addDataEntryByName(name string, data common.Blob) {
	// TODO: check for duplicate name
	d.strings = append(d.strings, name)
	d.nameEntries = append(d.nameEntries, &resourceDirectoryEntry{
		name: &name,
		dataEntry: &resourceDataEntry{
			data: data,
		},
	})
	d.sort()
}

func (d *resourceDirectory) addDataEntryByID(id int, data common.Blob) {
	d.idEntries = append(d.idEntries, &resourceDirectoryEntry{
		id: &id,
		dataEntry: &resourceDataEntry{
			data: data,
		},
	})
	d.sort()
}

func (d *resourceDirectory) addSubdirectoryByName(name string, characteristics uint32) {
	// TODO: check for duplicate name
	d.strings = append(d.strings, name)
	d.nameEntries = append(d.nameEntries, &resourceDirectoryEntry{
		name: &name,
		subdirectory: &resourceDirectory{
			characteristics: characteristics,
		},
	})
	d.sort()
}

func (d *resourceDirectory) addSubdirectoryByID(id int, characteristics uint32) {
	d.idEntries = append(d.idEntries, &resourceDirectoryEntry{
		id: &id,
		subdirectory: &resourceDirectory{
			characteristics: characteristics,
		},
	})
	d.sort()
}

func (d *resourceDirectory) sort() {
	sort.SliceStable(d.nameEntries, func(i, j int) bool {
		return *d.nameEntries[i].name < *d.nameEntries[j].name
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
		for _, e := range dir.entries() {
			if e.subdirectory != nil {
				if err := _walk(e.subdirectory); err != nil {
					return err
				}
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

type rawResourceDirectoryEntry struct {
	NameOffsetOrIntegerID               uint32
	DataEntryOffsetOrSubdirectoryOffset uint32
}

type resourceDirectoryEntry struct {
	offset       uint32
	name         *string
	id           *int
	dataEntry    *resourceDataEntry
	subdirectory *resourceDirectory
}

type resourceString struct {
	offset uint32
	length uint16
	s      string
}
