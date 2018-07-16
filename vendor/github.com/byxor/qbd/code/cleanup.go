package code

import (
	. "github.com/byxor/qbd/tokens"
)

func clean(tokens []Token) []Token {
	cleanTokens := make([]Token, len(tokens))

	cleanTokenCount := 0
	lastToken := Token{Invalid, nil}

	for _, token := range tokens {
		if !(token.Type == EndOfLine && lastToken.Type == EndOfLine) {
			cleanTokens[cleanTokenCount] = token
			cleanTokenCount++
		}
		lastToken = token
	}

	return cleanTokens[:cleanTokenCount]
}
