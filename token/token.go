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

	StartOfIf
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
	{isEndOfIf, EndOfIf},
	{isName, Name},
	{isInteger, Integer},
	{isStartOfFunction, StartOfFunction},
	{isEndOfFunction, EndOfFunction},
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

func isEndOfFile(bytes []byte) bool {
	return bytes[0] == 0x00
}

func isEndOfLine(bytes []byte) bool {
	return bytes[0] == 0x01
}

func isStartOfStruct(bytes []byte) bool {
	return bytes[0] == 0x03
}

func isEndOfStruct(bytes []byte) bool {
	return bytes[0] == 0x04
}

func isStartOfArray(bytes []byte) bool {
	return bytes[0] == 0x05
}

func isEndOfArray(bytes []byte) bool {
	return bytes[0] == 0x06
}

func isStartOfFunction(bytes []byte) bool {
	return bytes[0] == 0x23
}

func isEndOfFunction(bytes []byte) bool {
	return bytes[0] == 0x24
}

func isStartOfIf(bytes []byte) bool {
	return bytes[0] == 0x25
}

func isEndOfIf(bytes []byte) bool {
	return bytes[0] == 0x26
}

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
