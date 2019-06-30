package rsrc

import "github.com/hallazzang/syso/pkg/common"

type rawDataEntry struct {
	DataRVA  uint32
	Size     uint32
	Codepage uint32
	Reserved uint32
}

// DataEntry is a header for Data.
type DataEntry struct {
	offset uint32
	data   *Data
}

// Data represents actual binary resource data in .rsrc section.
type Data struct {
	offset uint32
	common.Blob
}
