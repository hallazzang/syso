package coff

type rawCOFFSymbol struct {
	Name               uint64
	Value              uint32
	SectionNumber      uint16
	Type               uint16
	StorageClass       uint8
	NumberOfAuxSymbols uint8
}

type coffSymbol struct {
	offset           uint32
	name             *coffString
	val              int
	sectionNo        int
	typ              int
	storageCls       int
	auxiliarySymbols []*coffAuxiliarySymbol
}

type coffString struct {
	offset uint32
	string
}

// TODO: not implemented
type coffAuxiliarySymbol struct {
}
