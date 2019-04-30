package domain

type Checksum uint32

type ChecksumGenerator interface {
	Generate(identifier string) Checksum
}

type ChecksumEncoder interface {
	LittleEndian(checksum Checksum) []byte
}
