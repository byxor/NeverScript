package token

type Token int

const (
	EndOfFile Token = iota
	EndOfLine
	Name
	Invalid
)

func GetTokens(tokens chan Token, input []byte) {
	if len(input) == 0 {
		tokens <- Invalid
	} else {
		first := input[0]
		if first == 0x00 {
			tokens <- EndOfFile
		} else if first == 0x01 {
			tokens <- EndOfLine
		} else if (first == 0x16) && len(input) == 5 {
			tokens <- Name
		} else {
			tokens <- Invalid
		}
	}
}
