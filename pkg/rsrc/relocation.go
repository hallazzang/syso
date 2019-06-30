package rsrc

type Relocation struct {
	va uint32
}

func (r *Relocation) VirtualAddress() uint32 {
	return r.va
}
