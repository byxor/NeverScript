package token

type Token int

const (
	EndOfFile Token = iota
	EndOfLine
)

func GetTokens(tokens chan Token, input []byte) {
	first := input[0]
	if first == 0x00 {
		tokens <- EndOfFile
	} else {
		tokens <- EndOfLine
	}
}
