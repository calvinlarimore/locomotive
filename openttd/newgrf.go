package openttd

type md5 [16]byte

type newGRFIdentifier struct {
	ID   uint32
	MD5  md5
	Name string
}

type newGRFConfig struct {
	ID newGRFIdentifier
}
