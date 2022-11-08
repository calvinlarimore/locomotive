package openttd

const (
	gameVersion          = "12.2"
	gameVersionMajor     = 12
	gameVersionMinor     = 2
	newGRFRevision       = (uint32(gameVersionMajor+16) << 24) | (uint32(gameVersionMinor) << 20) | (0x01 << 19) | 28004
	networkNameLength    = uint(80)
	networkVersionLength = uint(33)
)
