package coff

type rawRelocation struct {
	VirtualAddress   uint32
	SymbolTableIndex uint32
	Type             uint16
}

type relocation struct {
	offset uint32
	Relocation
}

type Relocation interface {
	VirtualAddress() uint32
}
