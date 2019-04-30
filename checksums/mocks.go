package checksums

import (
	"github.com/byxor/NeverScript"
)

type mockService struct{}

func NewMockService() Service {
	return &mockService{}
}

func (this mockService) GenerateFrom(identifier string) NeverScript.Checksum {
	return NeverScript.NewEmptyChecksum()
}

func (this mockService) EncodeAsLittleEndian(checksum NeverScript.Checksum) []byte {
	return []byte{}
}
