package table

import (
	. "github.com/byxor/qbd/tokens"
)

type NameTable struct {
	table map[int]string
}

func GenerateUsing(tokens []Token) NameTable {
	return NameTable{
		map[int]string{
			0: "Hello",
			1: "Hi",
		},
	}
}

func (nt NameTable) Get(checksum int) string {
	return ""
}
