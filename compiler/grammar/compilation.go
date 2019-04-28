package grammar

import (
	"github.com/byxor/NeverScript/compiler/checksums"
)

func (self EndOfLine) Compile() []byte {
	return []byte{0x01}
}

func (self Value) Compile() []byte {
	return []byte{0x01, 0x02, 0x03, 0x04}
}

func (self Assignment) Compile() []byte {
	name := []byte{0x66, 0x6F, 0x6F, 0x00}

	checksumBytes := checksums.LittleEndian(checksums.Generate(string(name)))

	return []byte{
		0x016,
		checksumBytes[0], checksumBytes[1], checksumBytes[2], checksumBytes[3],
		0x07,
		0x17, 0x0A, 0x00, 0x00, 0x00,
	}
}
