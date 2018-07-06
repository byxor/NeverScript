package token

func GetTokens(tokens chan Token, bytes []byte) {
	if len(bytes) == 0 {
		close(tokens)
		return
	}
	
	for chunkSize := 1; chunkSize <= len(bytes); chunkSize++ {

		chunk := bytes[:chunkSize]

		for _, c := range constructors {
			if c.function(chunk) {
				tokens <- c.token
				GetTokens(tokens, bytes[chunkSize:])
				return
			}
		}
	}

	tokens <- Invalid
}

type Token int

const (
	EndOfFile Token = iota
	EndOfLine
	Assignment
	Subtraction
	Addition
	Division
	Multiplication
	EqualityCheck
	LessThanCheck
	LessThanOrEqualCheck
	GreaterThanCheck
	GreaterThanOrEqualCheck
	StartOfStruct
	EndOfStruct
	StartOfArray
	EndOfArray
	StartOfFunction
	EndOfFunction
	Return
	Break
	StartOfIf
	Else
	ElseIf
	EndOfIf
	Integer
	Float
	Name
	ChecksumTableEntry
	Invalid
)

// -----------------------------

type constructor struct {
	token    Token
	function func([]byte) bool
}

/* The constructor functions are checked in order.
 * The ordering is important! */
var constructors = []constructor{
	{EndOfFile, singleByte(0x00)},
	{EndOfLine, singleByte(0x01)},
	{StartOfStruct, singleByte(0x03)},
	{EndOfStruct, singleByte(0x04)},
	{StartOfArray, singleByte(0x05)},
	{EndOfArray, singleByte(0x06)},
	{Assignment, singleByte(0x07)},
	{EqualityCheck, singleByte(0x11)},
	{LessThanCheck, singleByte(0x12)},
	{LessThanOrEqualCheck, singleByte(0x13)},
	{GreaterThanCheck, singleByte(0x14)},
	{GreaterThanOrEqualCheck, singleByte(0x15)},
	{Subtraction, singleByte(0x0A)},
	{Addition, singleByte(0x0B)},
	{Division, singleByte(0x0C)},
	{Multiplication, singleByte(0x0D)},
	{Break, singleByte(0x22)},
	{StartOfFunction, singleByte(0x23)},
	{EndOfFunction, singleByte(0x24)},
	{StartOfIf, singleByte(0x25)},
	{Else, singleByte(0x26)},
	{ElseIf, singleByte(0x27)},
	{EndOfIf, singleByte(0x28)},
	{Return, singleByte(0x29)},
	{Name, isName},
	{Integer, isInteger},
	{Float, isFloat},
	{ChecksumTableEntry, isCheckSumTableEntry},
}

func isName(bytes []byte) bool {
	return singleByte(0x16)(bytes) && len(bytes) == 5
}

func isInteger(bytes []byte) bool {
	return singleByte(0x17)(bytes) && len(bytes) == 5
}

func isFloat(bytes []byte) bool {
	return singleByte(0x1A)(bytes) && len(bytes) == 5
}

func isCheckSumTableEntry(bytes []byte) bool {
	isLongEnough := len(bytes) > 6
	isNullTerminated := bytes[len(bytes)-1] == 0
	return singleByte(0x2B)(bytes) && isLongEnough && isNullTerminated
}

func singleByte(n byte) func([]byte) bool {
	return func(bytes []byte) bool {
		return bytes[0] == n
	}
}
