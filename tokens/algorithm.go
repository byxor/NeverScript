package tokens

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
)

func ExtractAll(tokenChannel chan Token, chunk []byte) {
	if len(chunk) == 0 {
		close(tokenChannel)
	} else {
		token, gotOne := searchForToken(chunk)
		if gotOne {
			color.Green(token.Type.String())
			color.White(hex.Dump(token.Chunk) + "\n")
			tokenChannel <- token
			nextChunk := chunk[len(token.Chunk):]
			ExtractAll(tokenChannel, nextChunk)
		} else {
			color.Yellow(fmt.Sprintf("Unrecognised chunk\n%s\n", hex.Dump(token.Chunk)))
			tokenChannel <- Token{Invalid, token.Chunk}
			close(tokenChannel)
		}
	}
}

func searchForToken(chunk []byte) (token Token, gotOne bool) {
	for subChunkSize := 1; subChunkSize <= len(chunk); subChunkSize++ {
		subChunk := chunk[:subChunkSize]

		tokenType, gotOne := checkIfChunkRepresentsToken(subChunk)
		if gotOne {
			return Token{tokenType, subChunk}, gotOne
		}
	}
	return Token{Invalid, chunk}, false
}

func checkIfChunkRepresentsToken(chunk []byte) (tokenType TokenType, gotOne bool) {
	var constructors []constructor

	if len(chunk) == 1 {
		constructors = singleByteConstructors
	} else {
		constructors = otherConstructors
	}

	for _, c := range constructors {
		if c.function(chunk) {
			return c.tokenType, true
		}
	}

	return Invalid, false
}

func readInt32(bytes []byte) int {
	return int(binary.LittleEndian.Uint32(bytes))
}
