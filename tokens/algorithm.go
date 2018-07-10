package tokens

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
)

func Extract(tokenChannel chan Token, chunk []byte) {
	if len(chunk) == 0 {
		close(tokenChannel)
		return
	}

	token, subChunk, gotOne := searchForToken(chunk)
	if gotOne {
		color.Green(token.String())
		color.White(hex.Dump(subChunk) + "\n")

		tokenChannel <- token
		nextChunk := chunk[len(subChunk):]
		Extract(tokenChannel, nextChunk)
	} else {
		color.Yellow(fmt.Sprintf("Unrecognised chunk\n%s\n", hex.Dump(subChunk)))
		tokenChannel <- Invalid
		close(tokenChannel)
	}
}

func searchForToken(chunk []byte) (token Token, subChunk []byte, gotOne bool) {
	for subChunkSize := 1; subChunkSize <= len(chunk); subChunkSize++ {
		subChunk := chunk[:subChunkSize]

		token, gotOne := checkIfChunkRepresentsToken(subChunk)
		if gotOne {
			return token, subChunk, gotOne
		}
	}
	return Invalid, chunk, false
}

func checkIfChunkRepresentsToken(chunk []byte) (token Token, gotOne bool) {
	var constructors []constructor

	if len(chunk) == 1 {
		constructors = singleByteConstructors
	} else {
		constructors = otherConstructors
	}

	for _, c := range constructors {
		if c.function(chunk) {
			return c.token, true
		}
	}

	return Invalid, false
}

func readInt32(bytes []byte) int {
	return int(binary.LittleEndian.Uint32(bytes))
}
