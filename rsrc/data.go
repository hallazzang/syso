package rsrc

import "github.com/hallazzang/syso/internal/common"

type resourceDataEntry struct {
	DataRVA  uint32
	Size     uint32
	Codepage uint32
	Reserved uint32
}

type dataEntry struct {
	offset uint32
	data   common.Blob
}
