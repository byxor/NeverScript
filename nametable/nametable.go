package nametable

import (
	"encoding/hex"
	"fmt"
	. "github.com/byxor/qbd/tokens"
)

type NameTable struct {
	entries internalTable
}

type internalTable map[string]string

func GenerateUsing(tokens []Token) NameTable {
	var table = make(internalTable)

	for _, token := range tokens {
		if token.Type == NameTableEntry {
			checksum := hex.EncodeToString(reverse(token.Chunk[1:5]))
			fmt.Println(checksum)
			name := string(token.Chunk[5 : len(token.Chunk)-1])
			table[checksum] = name
		}
	}

	return NameTable{table}
}

func (nt NameTable) Get(checksum string) string {
	fmt.Println(nt.entries)
	if name, ok := nt.entries[checksum]; ok {
		return name
	} else {
		return "&" + checksum
	}
}

func reverse(chunk []byte) []byte {
	length := len(chunk)
	reversed := make([]byte, length)
	for i, v := range chunk {
		reversed[length-1-i] = v
	}
	return reversed
}
