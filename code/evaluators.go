package code

import (
	"encoding/hex"
	. "github.com/byxor/qbd/tokens"
	"strconv"
)

type evaluator func(*stateHolder) string

var evaluators = map[TokenType]evaluator{
	EndOfFile:         basicString(""),
	EndOfLine:         basicString("\n; "),
	StartOfArray:      evaluateStartOfArray,
	EndOfArray:        evaluateEndOfArray,
	StartOfExpression: evaluateStartOfExpression,
	EndOfExpression:   evaluateEndOfExpression,
	Assignment:        basicString(" = "),
	Addition:          basicString(" + "),
	Subtraction:       basicString(" - "),
	Multiplication:    basicString(" * "),
	Division:          basicString(" / "),
	LocalReference:    evaluateLocalReference,
	Integer:           evaluateInteger,
	Name:              evaluateName,
	NameTableEntry:    basicString(""),
}

func evaluateStartOfExpression(state *stateHolder) string {
	state.neverAddWhitespace = true
	return basicString("(")(state)
}

func evaluateEndOfExpression(state *stateHolder) string {
	state.neverAddWhitespace = true
	return basicString(")")(state)
}

func evaluateLocalReference(state *stateHolder) string {
	state.neverAddWhitespace = true
	return basicString("$")(state)
}

func evaluateStartOfArray(state *stateHolder) string {
	state.arrayDepth++
	return basicString("[")(state)
}

func evaluateEndOfArray(state *stateHolder) string {
	state.arrayDepth--
	return basicString("]")(state)
}

func evaluateInteger(state *stateHolder) string {
	return strconv.Itoa(ReadInt32(state.token.Chunk[1:]))
}

func evaluateName(state *stateHolder) string {
	checksum := hex.EncodeToString(state.token.Chunk[1:])
	return state.names.Get(checksum)
}

func basicString(s string) evaluator {
	return func(*stateHolder) string {
		return s
	}
}
