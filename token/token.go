package token

type Token int

const (
	EndOfFile Token = iota
	EndOfLine

	StartOfStruct
	EndOfStruct

	StartOfArray
	EndOfArray

	StartOfFunction
	EndOfFunction
	Return

	StartOfIf
	Else
	ElseIf
	EndOfIf

	Name
	Integer

	Invalid
	None
)

/* The constructor functions are checked in order.
 * The ordering is important! */
var constructors = []constructor{
	{isEndOfFile, EndOfFile},
	{isEndOfLine, EndOfLine},
	{isStartOfStruct, StartOfStruct},
	{isEndOfStruct, EndOfStruct},
	{isStartOfArray, StartOfArray},
	{isEndOfArray, EndOfArray},
	{isStartOfIf, StartOfIf},
	{isElse, Else},
	{isElseIf, ElseIf},
	{isEndOfIf, EndOfIf},
	{isName, Name},
	{isInteger, Integer},
	{isStartOfFunction, StartOfFunction},
	{isEndOfFunction, EndOfFunction},
	{isReturn, Return},
}

func GetTokens(tokens chan Token, bytes []byte) {
	if len(bytes) == 0 {
		tokens <- None
		return
	}

	for _, c := range constructors {
		if c.function(bytes) {
			tokens <- c.token
			return
		}
	}

	tokens <- Invalid
}

// -----------------------------

type constructor struct {
	function func([]byte) bool
	token    Token
}

var isEndOfFile = singleByte(0x00)
var isEndOfLine = singleByte(0x01)
var isStartOfStruct = singleByte(0x03)
var isEndOfStruct = singleByte(0x04)
var isStartOfArray = singleByte(0x05)
var isEndOfArray = singleByte(0x06)
var isStartOfFunction = singleByte(0x23)
var isEndOfFunction = singleByte(0x24)
var isStartOfIf = singleByte(0x25)
var isElse = singleByte(0x26)
var isElseIf = singleByte(0x27)
var isEndOfIf = singleByte(0x28)
var isReturn = singleByte(0x29)

func isName(bytes []byte) bool {
	hasPrefix := bytes[0] == 0x16
	longEnough := len(bytes) == 5
	return hasPrefix && longEnough
}

func isInteger(bytes []byte) bool {
	hasPrefix := bytes[0] == 0x17
	longEnough := len(bytes) == 5
	return hasPrefix && longEnough
}

func singleByte(n byte) func([]byte) bool {
	return func(bytes []byte) bool {
		return bytes[0] == n
	}
}
