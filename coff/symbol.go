package coff

type rawSymbol struct {
	Name               [8]byte
	Value              uint32
	SectionNumber      uint16
	Type               uint16
	StorageClass       uint8
	NumberOfAuxSymbols uint8
}
