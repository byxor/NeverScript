package token

type Token int

const (
	EndOfFile Token = iota
	EndOfLine
	Name
	Integer
	Invalid
)

func GetTokens(tokens chan Token, bytes []byte) {
	if len(bytes) == 0 {
		tokens <- Invalid
	} else {
		if isEndOfFile(bytes) {
			tokens <- EndOfFile
		} else if isEndOfLine(bytes) {
			tokens <- EndOfLine
		} else if isName(bytes) {
			tokens <- Name
		} else if isInteger(bytes) {
			tokens <- Integer
		} else {
			tokens <- Invalid
		}
	}
}

func isEndOfFile(bytes []byte) bool {
	return bytes[0] == 0x00
}

func isEndOfLine(bytes []byte) bool {
	return bytes[0] == 0x01
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
