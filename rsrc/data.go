package rsrc

type rawResourceDataEntry struct {
	DataRVA  uint32
	Size     uint32
	Codepage uint32
	Reserved uint32
}

type resourceDataEntry struct {
	offset uint32
	data   Blob
}
