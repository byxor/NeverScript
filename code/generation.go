package code

import (
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
	return ""
}

func evaluateInteger(chunk []byte) string {
	return strconv.Itoa(ReadInt32(chunk))
}
