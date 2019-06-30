package rsrc

import "github.com/hallazzang/syso/pkg/common"

type rawDataEntry struct {
	DataRVA  uint32
	Size     uint32
	Codepage uint32
	Reserved uint32
}

type DataEntry struct {
	offset uint32
	data   *Data
}

type Data struct {
	offset uint32
	common.Blob
}
