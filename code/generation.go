package code

import (
	"encoding/hex"
	"github.com/byxor/qbd/nametable"
	. "github.com/byxor/qbd/tokens"
	"strconv"
)

func GenerateUsing(tokens []Token) string {
	return generateUsing(clean(tokens), nametable.BuildFrom(tokens))
}

func clean(tokens []Token) []Token {
	cleanTokens := make([]Token, len(tokens))

	cleanTokenCount := 0
	lastToken := Token{Invalid, nil}

	for _, token := range tokens {
		if !(token.Type == EndOfLine && lastToken.Type == EndOfLine) {
			cleanTokens[cleanTokenCount] = token
			cleanTokenCount++
		}
		lastToken = token
	}

	return cleanTokens[:cleanTokenCount]
}

func generateUsing(tokens []Token, nameTable nametable.NameTable) string {
	if len(tokens) == 0 {
		return ""
	}
	evaluator := evaluators[tokens[0].Type]
	result := evaluator(tokens[0].Chunk, nameTable)
	return result + generateUsing(tokens[1:], nameTable)
}

type evaluator func([]byte, nametable.NameTable) string

var evaluators = map[TokenType]evaluator{
	EndOfFile:      basicString(""),
	EndOfLine:      basicString(";"),
	Addition:       basicString(" + "),
	Integer:        evaluateInteger,
	Name:           evaluateName,
	NameTableEntry: basicString(""),
}

func evaluateInteger(chunk []byte, nameTable nametable.NameTable) string {
	return strconv.Itoa(ReadInt32(chunk[1:]))
}

func evaluateName(chunk []byte, nameTable nametable.NameTable) string {
	checksum := hex.EncodeToString(chunk[1:])
	return nameTable.Get(checksum)
}

func basicString(s string) evaluator {
	return func([]byte, nametable.NameTable) string {
		return s
	}
}
