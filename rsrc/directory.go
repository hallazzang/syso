package rsrc

import (
	"sort"

	"github.com/hallazzang/syso/internal/common"
)

type resourceDirectory struct {
	Characteristics     uint32
	TimeDateStamp       uint32
	MajorVersion        uint16
	MinorVersion        uint16
	NumberOfNameEntries uint16
	NumberOfIDEntries   uint16
}

type directory struct {
	offset          uint32
	characteristics uint32
	nameEntries     []*directoryEntry
	idEntries       []*directoryEntry
	names           []string
}

func (d *directory) addDataEntryByName(name string, data common.Blob) {
	// TODO: check for duplicate name
	d.nameEntries = append(d.nameEntries, &directoryEntry{
		name: &name,
		dataEntry: &dataEntry{
			data: data,
		},
	})
	d.sort()
}

func (d *directory) addDataEntryByID(id int, data common.Blob) {
	d.idEntries = append(d.idEntries, &directoryEntry{
		id: &id,
		dataEntry: &dataEntry{
			data: data,
		},
	})
	d.sort()
}

func (d *directory) addSubdirectoryByName(name string, characteristics uint32) {
	// TODO: check for duplicate name
	d.nameEntries = append(d.nameEntries, &directoryEntry{
		name: &name,
		subdirectory: &directory{
			characteristics: characteristics,
		},
	})
	d.sort()
}

func (d *directory) addSubdirectoryByID(id int, characteristics uint32) {
	d.idEntries = append(d.idEntries, &directoryEntry{
		id: &id,
		subdirectory: &directory{
			characteristics: characteristics,
		},
	})
	d.sort()
}

func (d *directory) sort() {
	sort.SliceStable(d.nameEntries, func(i, j int) bool {
		return *d.nameEntries[i].name < *d.nameEntries[j].name
	})
	sort.SliceStable(d.idEntries, func(i, j int) bool {
		return *d.idEntries[i].id < *d.idEntries[j].id
	})
}

func (d *directory) walk(cb func(*directory) error) error {
	var _walk func(*directory) error
	_walk = func(dir *directory) error {
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

func (d *directory) entries() []*directoryEntry {
	return append(append([]*directoryEntry{}, d.nameEntries...), d.idEntries...)
}

type resourceDirectoryEntry struct {
	NameOffsetOrIntegerID               uint32
	DataEntryOffsetOrSubdirectoryOffset uint32
}

type directoryEntry struct {
	offset       uint32
	name         *string
	id           *int
	dataEntry    *dataEntry
	subdirectory *directory
}

type str struct {
	offset uint32
	length uint16
	s      string
}
