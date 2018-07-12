package code

import (
	"encoding/hex"
	"github.com/byxor/qbd/nametable"
	. "github.com/byxor/qbd/tokens"
	"strconv"
)

type evaluator func([]byte, nametable.NameTable) string

var evaluators = map[TokenType]evaluator{
	EndOfFile:      basicString(""),
	EndOfLine:      basicString("; "),
	Addition:       basicString(" + "),
	Subtraction:    basicString(" - "),
	Multiplication: basicString(" * "),
	Division:       basicString(" / "),
	LocalReference: basicString("$"),
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
