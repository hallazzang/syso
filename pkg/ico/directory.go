package ico

// ICONDIR
type directory struct {
	Reserved uint16
	Type     uint16
	Count    uint16
}

// ICONDIRENTRY
type directoryEntry struct {
	Width       uint8
	Height      uint8
	ColorCount  uint8
	Reserved    uint8
	Planes      uint16
	BitCount    uint16
	BytesInRes  uint32
	ImageOffset uint32
}

// GRPICONDIR
type groupDirectory directory

// GRPICONDIRENTRY
type groupDirectoryEntry struct {
	Width      uint8
	Height     uint8
	ColorCount uint8
	Reserved   uint8
	Planes     uint16
	BitCount   uint16
	BytesInRes uint32
	ID         uint16
}
