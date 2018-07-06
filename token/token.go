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
	{EndOfFile, newSingleByteChecker(0x00)},
	{EndOfLine, newSingleByteChecker(0x01)},
	{StartOfStruct, newSingleByteChecker(0x03)},
	{EndOfStruct, newSingleByteChecker(0x04)},
	{StartOfArray, newSingleByteChecker(0x05)},
	{EndOfArray, newSingleByteChecker(0x06)},
	{Assignment, newSingleByteChecker(0x07)},
	{EqualityCheck, newSingleByteChecker(0x11)},
	{LessThanCheck, newSingleByteChecker(0x12)},
	{LessThanOrEqualCheck, newSingleByteChecker(0x13)},
	{GreaterThanCheck, newSingleByteChecker(0x14)},
	{GreaterThanOrEqualCheck, newSingleByteChecker(0x15)},
	{Subtraction, newSingleByteChecker(0x0A)},
	{Addition, newSingleByteChecker(0x0B)},
	{Division, newSingleByteChecker(0x0C)},
	{Multiplication, newSingleByteChecker(0x0D)},
	{Break, newSingleByteChecker(0x22)},
	{StartOfFunction, newSingleByteChecker(0x23)},
	{EndOfFunction, newSingleByteChecker(0x24)},
	{StartOfIf, newSingleByteChecker(0x25)},
	{Else, newSingleByteChecker(0x26)},
	{ElseIf, newSingleByteChecker(0x27)},
	{EndOfIf, newSingleByteChecker(0x28)},
	{Return, newSingleByteChecker(0x29)},
	{Name, isName},
	{Integer, isInteger},
	{Float, isFloat},
	{ChecksumTableEntry, isCheckSumTableEntry},
}

func isName(bytes []byte) bool {
	return newSingleByteChecker(0x16)(bytes) && len(bytes) == 5
}

func isInteger(bytes []byte) bool {
	return newSingleByteChecker(0x17)(bytes) && len(bytes) == 5
}

func isFloat(bytes []byte) bool {
	return newSingleByteChecker(0x1A)(bytes) && len(bytes) == 5
}

func isCheckSumTableEntry(bytes []byte) bool {
	isLongEnough := len(bytes) > 6
	isNullTerminated := bytes[len(bytes)-1] == 0
	return newSingleByteChecker(0x2B)(bytes) && isLongEnough && isNullTerminated
}

func newSingleByteChecker(n byte) func([]byte) bool {
	return func(bytes []byte) bool {
		return bytes[0] == n
	}
}
