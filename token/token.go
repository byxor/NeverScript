package token

type Token int

const (
	EndOfFile Token = iota
	EndOfLine
	StartOfStruct
	EndOfStruct
	StartOfArray
	EndOfArray
	Name
	Integer

	Invalid
	None
)

var constructors = []constructor{
	{isEndOfFile, EndOfFile},
	{isEndOfLine, EndOfLine},
	{isStartOfStruct, StartOfStruct},
	{isEndOfStruct, EndOfStruct},
	{isStartOfArray, StartOfArray},
	{isEndOfArray, EndOfArray},
	{isName, Name},
	{isInteger, Integer},
}

func GetTokens(tokens chan Token, bytes []byte) {
	if len(bytes) == 0 {
		tokens <- None
	} else {
		for _, c := range constructors {
			if c.function(bytes) {
				tokens <- c.token
				break
			}
		}
		tokens <- Invalid
	}
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
