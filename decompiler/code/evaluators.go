package code

import (
	"encoding/hex"
	"fmt"
	"strconv"

	. "github.com/byxor/NeverScript/decompiler/tokens"
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
	Float:             evaluateFloat,
	Pair:              evaluatePair,
	Vector:            evaluateVector,
	Name:              evaluateName,
	NameTableEntry:    basicString(""),
	String:            evaluateString,
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

func evaluateFloat(state *stateHolder) string {
	return formatFloat(ReadFloat32(state.token.Chunk[1:]))
}

func evaluatePair(state *stateHolder) string {
	firstValue := ReadFloat32(state.token.Chunk[1:5])
	secondValue := ReadFloat32(state.token.Chunk[5:])
	return fmt.Sprintf(
		"vec2<%s, %s>",
		formatFloat(firstValue),
		formatFloat(secondValue),
	)
}

func evaluateVector(state *stateHolder) string {
	firstValue := ReadFloat32(state.token.Chunk[1:5])
	secondValue := ReadFloat32(state.token.Chunk[5:9])
	thirdValue := ReadFloat32(state.token.Chunk[9:])
	return fmt.Sprintf(
		"vec3<%s, %s, %s>",
		formatFloat(firstValue),
		formatFloat(secondValue),
		formatFloat(thirdValue),
	)
}

func evaluateName(state *stateHolder) string {
	checksum := hex.EncodeToString(state.token.Chunk[1:])
	return state.names.Get(checksum)
}

func evaluateString(state *stateHolder) string {
	chunk := state.token.Chunk
	return "\"" + string(chunk[5:len(chunk)-1]) + "\""
}

func basicString(s string) evaluator {
	return func(*stateHolder) string {
		return s
	}
}

func formatFloat(f float32) string {
	return trimUnnecessaryDigits(fmt.Sprintf("%.7f", f))
}

func trimUnnecessaryDigits(floatString string) string {
	end := 0
	for i := len(floatString) - 1; i >= 0; i-- {
		if floatString[i-1] == '.' {
			end = i
			break
		}
		if floatString[i] != '0' {
			end = i
			break
		}
	}
	return floatString[:end+1]
}
