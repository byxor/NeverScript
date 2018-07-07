package tokens

import (
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
)

func Extract(tokenChannel chan Token, bytes []byte) {
	if len(bytes) == 0 {
		close(tokenChannel)
		return
	}

	var chunk []byte

	for chunkSize := 1; chunkSize <= len(bytes); chunkSize++ {
		chunk = bytes[:chunkSize]

		if chunkSize >= 50 {
			color.Yellow("Giving up...")
			break
		}

		for _, c := range constructors {
			if c.function(chunk) {

				color.Green(c.token.String())
				color.White(hex.Dump(chunk))
				color.White("")

				tokenChannel <- c.token
				Extract(tokenChannel, bytes[chunkSize:])
				return
			}
		}
	}

	color.Yellow(fmt.Sprintf("Unrecognised chunk\n%s\n", hex.Dump(chunk)))

	tokenChannel <- Invalid
	close(tokenChannel)
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
	{Break, requirePrefix(0x22)},
	{StartOfFunction, requirePrefix(0x23)},
	{EndOfFunction, requirePrefix(0x24)},
	{StartOfIf, requirePrefix(0x25)},
	{Else, requirePrefix(0x26)},
	{ElseIf, requirePrefix(0x27)},
	{EndOfIf, requirePrefix(0x28)},
	{Return, requirePrefix(0x29)},
	{WeirdThing, requirePrefix(0x42)},
	{Name, requirePrefixAndLength(0x16, 5)},
	{Integer, requirePrefixAndLength(0x17, 5)},
	{Float, requirePrefixAndLength(0x1A, 5)},
	{ChecksumTableEntry, isCheckSumTableEntry},
}

func isCheckSumTableEntry(bytes []byte) bool {
	isLongEnough := len(bytes) > 6
	isNullTerminated := bytes[len(bytes)-1] == 0
	return requirePrefix(0x2B)(bytes) && isLongEnough && isNullTerminated
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
