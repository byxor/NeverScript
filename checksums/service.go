package checksums

import (
	"github.com/byxor/NeverScript"
)

type Service interface {
	Generate(identifier string) NeverScript.Checksum
	EncodeAsLittleEndian(checksum NeverScript.Checksum) []byte
}
