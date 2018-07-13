package code

import (
	"encoding/hex"
	"github.com/byxor/qbd/nametable"
	. "github.com/byxor/qbd/tokens"
	"strconv"
)

type evaluator func(stateMap) string

var evaluators = map[TokenType]evaluator{
	EndOfFile:      basicString(""),
	EndOfLine:      basicString("; "),
	StartOfArray:   evaluateStartOfArray,
	EndOfArray:     basicString("]"),
	Assignment:     basicString(" = "),
	Addition:       basicString(" + "),
	Subtraction:    basicString(" - "),
	Multiplication: basicString(" * "),
	Division:       basicString(" / "),
	LocalReference: basicString("$"),
	Integer:        evaluateInteger,
	Name:           evaluateName,
	NameTableEntry: basicString(""),
}

func evaluateStartOfArray(state stateMap) string {
	state["inArray"] = 1
	return basicString("[")(state)
}

func evaluateInteger(state stateMap) string {
	token := state["token"].(Token)
	return strconv.Itoa(ReadInt32(token.Chunk[1:]))
}

func evaluateName(state stateMap) string {
	token := state["token"].(Token)
	checksum := hex.EncodeToString(token.Chunk[1:])
	return state["names"].(nametable.NameTable).Get(checksum)
}

func basicString(s string) evaluator {
	return func(stateMap) string {
		return s
	}
}
