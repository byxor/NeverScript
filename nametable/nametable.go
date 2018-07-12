package nametable

import (
	"encoding/hex"
	. "github.com/byxor/qbd/tokens"
)

func GenerateUsing(tokens []Token) NameTable {
	var entries = make(internalTable)
	for _, token := range tokens {
		if token.Type == NameTableEntry {
			storeEntry(entries, token)
		}
	}
	return NameTable{entries}
}

func storeEntry(entries internalTable, token Token) {
	checksum := hex.EncodeToString(token.Chunk[1:5])
	name := string(token.Chunk[5 : len(token.Chunk)-1])
	entries[checksum] = name
}

func (nt NameTable) Get(checksum string) string {
	if name, ok := nt.entries[checksum]; ok {
		return name
	} else {
		return "&" + checksum
	}
}
