package code

import (
	"encoding/hex"
	"github.com/byxor/qbd/nametable"
	. "github.com/byxor/qbd/tokens"
	"strconv"
)

func GenerateUsing(tokens []Token) string {
	if len(tokens) == 0 {
		return ""
	}
	if tokens[0].Type == EndOfLine {
		return ";"
	}
	if tokens[0].Type == Integer {
		return evaluateInteger(tokens[0].Chunk[1:])
	}
	if tokens[0].Type == Name {
		checksum := hex.EncodeToString(tokens[0].Chunk[1:])
		return nametable.BuildFrom([]Token{}).Get(checksum)
	}
	return ""
}

func evaluateInteger(chunk []byte) string {
	return strconv.Itoa(ReadInt32(chunk))
}
