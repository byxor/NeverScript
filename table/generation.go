package table

import (
	"encoding/binary"
	"encoding/hex"
	. "github.com/byxor/qbd/tokens"
)

type NameTable struct {
	entries internalTable
}

type internalTable map[int]string

func GenerateUsing(tokens []Token) NameTable {
	var table = make(internalTable)

	for _, token := range tokens {
		if token.Type == ChecksumTableEntry {
			checksum := ReadInt32(token.Chunk[1:5])
			name := string(token.Chunk[5 : len(token.Chunk)-1])
			table[checksum] = name
		}
	}

	return NameTable{table}
}

func (nt NameTable) Get(checksum int) string {
	if name, ok := nt.entries[checksum]; ok {
		return name
	} else {
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, uint32(checksum))
		return "&" + hex.EncodeToString(bytes)
	}
}
