package tokens

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
)

func Extract(tokenChannel chan Token, chunk []byte) {
	if len(chunk) == 0 {
		close(tokenChannel)
		return
	}

	token, subChunk, gotOne := searchForToken(chunk)
	if gotOne {
		color.Green(token.String())
		color.White(hex.Dump(subChunk) + "\n")

		tokenChannel <- token
		nextChunk := chunk[len(subChunk):]
		Extract(tokenChannel, nextChunk)
	} else {
		color.Yellow(fmt.Sprintf("Unrecognised chunk\n%s\n", hex.Dump(subChunk)))
		tokenChannel <- Invalid
		close(tokenChannel)
	}
}

func searchForToken(chunk []byte) (token Token, subChunk []byte, gotOne bool) {
	for subChunkSize := 1; subChunkSize <= len(chunk); subChunkSize++ {
		if subChunkSize >= 100 {
			color.Yellow("Giving up...")
			break
		}

		subChunk := chunk[:subChunkSize]

		token, gotOne := checkIfChunkRepresentsToken(subChunk)
		if gotOne {
			return token, subChunk, gotOne
		}
	}

	return Invalid, chunk, false
}

func checkIfChunkRepresentsToken(chunk []byte) (token Token, gotOne bool) {
	for _, c := range constructors {
		if c.function(chunk) {
			return c.token, true
		}
	}
	return Invalid, false
}

type constructor struct {
	token    Token
	function func([]byte) bool
}

/* The constructor functions are checked in order.
 * The ordering is important! */
var constructors = []constructor{
	{EndOfFile, requirePrefix(0x00)},
	{EndOfLine, requirePrefix(0x01)},
	{StartOfStruct, requirePrefix(0x03)},
	{EndOfStruct, requirePrefix(0x04)},
	{StartOfArray, requirePrefix(0x05)},
	{EndOfArray, requirePrefix(0x06)},
	{Assignment, requirePrefix(0x07)},
	{Comma, requirePrefix(0x09)},
	{StartOfExpression, requirePrefix(0x0E)},
	{EndOfExpression, requirePrefix(0x0F)},
	{LocalReference, requirePrefix(0x2D)},
	{EqualityCheck, requirePrefix(0x11)},
	{LessThanCheck, requirePrefix(0x12)},
	{LessThanOrEqualCheck, requirePrefix(0x13)},
	{GreaterThanCheck, requirePrefix(0x14)},
	{GreaterThanOrEqualCheck, requirePrefix(0x15)},
	{Subtraction, requirePrefix(0x0A)},
	{Addition, requirePrefix(0x0B)},
	{Division, requirePrefix(0x0C)},
	{Multiplication, requirePrefix(0x0D)},
	{While, requirePrefix(0x20)},
	{Repeat, requirePrefix(0x21)},
	{Break, requirePrefix(0x22)},
	{StartOfFunction, requirePrefix(0x23)},
	{EndOfFunction, requirePrefix(0x24)},
	{StartOfIf, requirePrefix(0x25)},
	{Else, requirePrefix(0x26)},
	{ElseIf, requirePrefix(0x27)},
	{EndOfIf, requirePrefix(0x28)},
	{Return, requirePrefix(0x29)},
	{AllLocalReferences, requirePrefix(0x2C)},
	{And, requirePrefix(0x33)},
	{Not, requirePrefix(0x39)},
	{StartOfSwitch, requirePrefix(0x3C)},
	{EndOfSwitch, requirePrefix(0x3D)},
	{SwitchCase, requirePrefix(0x3E)},
	{DefaultSwitchCase, requirePrefix(0x3F)},
	{NamespaceAccess, requirePrefix(0x42)},
	{Name, requirePrefixAndLength(0x16, 5)},
	{Integer, requirePrefixAndLength(0x17, 5)},
	{Float, requirePrefixAndLength(0x1A, 5)},
	{Pair, requirePrefixAndLength(0x1F, 9)},
	{LongJump, requirePrefixAndLength(0x2E, 5)},
	{OptimisedIf, requirePrefixAndLength(0x47, 3)},
	{OptimisedElse, requirePrefixAndLength(0x48, 3)},
	{ShortJump, requirePrefixAndLength(0x49, 3)},
	{ChecksumTableEntry, isCheckSumTableEntry},
	{String, isString},
	{LocalString, isLocalString},
	// {ExecuteRandomBlock, isExecuteRandomBlock},
}

func requirePrefixAndLength(prefix byte, length int) func([]byte) bool {
	return func(bytes []byte) bool {
		return requirePrefix(prefix)(bytes) && len(bytes) == length
	}
}

func requirePrefix(n byte) func([]byte) bool {
	return func(bytes []byte) bool {
		return bytes[0] == n
	}
}

func isCheckSumTableEntry(bytes []byte) bool {
	isLongEnough := len(bytes) > 6
	isNullTerminated := bytes[len(bytes)-1] == 0
	return requirePrefix(0x2B)(bytes) && isLongEnough && isNullTerminated
}

func isLocalString(bytes []byte) bool {
	return hasStringComponent(bytes) && requirePrefix(0x1C)(bytes)
}

func isString(bytes []byte) bool {
	return hasStringComponent(bytes) && requirePrefix(0x1B)(bytes)
}

func hasStringComponent(bytes []byte) bool {
	const headerLength = 5
	length := len(bytes)
	if length < headerLength {
		return false
	}
	stringLength := int(binary.LittleEndian.Uint32(bytes[1:headerLength]))
	return length == headerLength+stringLength
}

// func isExecuteRandomBlock(chunk []byte) bool {
// 	prefixLength := 1

// 	numberOfBlocks := int(binary.LittleEndian.Uint32(chunk[prefixLength : prefixLength+4]))

// 	weightSectionLength := 2 * numberOfBlocks
// 	offsetSectionLength := 4 * numberOfBlocks
// 	headerLength := prefixLength + weightSectionLength + offsetSectionLength

// 	firstOffsetBytes := chunk[headerLength-offsetSectionLength : headerLength-offsetSectionLength+4]
// 	firstOffset := binary.LittleEndian.Uint32(firstOffsetBytes)
// 	firstCodeBlock := chunk[firstOffset:]

// 	var nextChunk := chunk
// 	for {
// 		token, subChunk, gotOne := searchForToken(nextChunk)

// 		if token == LongJump

// 	}

// 	expectedLength := 1 + weightSectionLength

// 	hasPrefix := requirePrefix(0x2F)(bytes)
// }
