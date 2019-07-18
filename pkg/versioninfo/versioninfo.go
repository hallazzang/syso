package versioninfo

type VersionInfo struct {
	fixedFileInfo  fixedFileInfo
	stringFileInfo *stringFileInfo
	varFileInfo    *varFileInfo
}

type fixedFileInfo struct {
	// structVersion uint32 // TODO: do we need it?
	fileVersion    uint64
	productVersion uint64
	fileFlagsMask  uint32
	fileFlags      uint32
	fileOS         uint32
	fileType       uint32
	fileSubtype    uint32
	fileDate       uint64
}

type rawFixedFileInfo struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

type stringFileInfo struct {
	stringTables []stringTable
}

type stringTable struct {
	language uint16
	codepage uint16
	strings  []_string
}

type _string struct {
	key   string
	value string
}

type varFileInfo struct {
	translations []translation
}

type translation struct {
	language uint16
	codepage uint16
}
