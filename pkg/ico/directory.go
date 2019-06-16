package ico

// Directory represents an icon directory(ICONDIR).
type Directory struct {
	Reserved uint16
	Type     uint16
	Count    uint16
}

// DirectoryEntry represents an icon directory entry(ICONDIRENTRY).
type DirectoryEntry struct {
	Width       uint8
	Height      uint8
	ColorCount  uint8
	Reserved    uint8
	Planes      uint16
	BitCount    uint16
	BytesInRes  uint32
	ImageOffset uint32
}

// GroupDirectory represents an icon group directory(GRPICONDIR).
type GroupDirectory Directory

// GroupDirectoryEntry represents an icon group directory entry(GRPICONDIRENTRY).
type GroupDirectoryEntry struct {
	Width      uint8
	Height     uint8
	ColorCount uint8
	Reserved   uint8
	Planes     uint16
	BitCount   uint16
	BytesInRes uint32
	ID         uint16
}
