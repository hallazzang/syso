package rsrc

// Relocation represents a relocation in .rsrc section.
type Relocation struct {
	va uint32
}

// VirtualAddress returns a virtual address where the relocation
// should be applied to.
func (r *Relocation) VirtualAddress() uint32 {
	return r.va
}
