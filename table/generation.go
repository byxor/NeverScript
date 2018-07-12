package table

import (
	"encoding/binary"
	"encoding/hex"
	. "github.com/byxor/qbd/tokens"
)

type NameTable struct {
	table internalTable
}

type internalTable map[int]string

func GenerateUsing(tokens []Token) NameTable {
	return NameTable{
		internalTable{
			0: "Hello",
			1: "Hi",
		},
	}
}

func (nt NameTable) Get(checksum int) string {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(checksum))
	return "&" + hex.EncodeToString(bytes)
}
