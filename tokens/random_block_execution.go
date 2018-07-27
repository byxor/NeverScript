package tokens

func isExecuteRandomBlock(chunk []byte) bool {

	prefix := byte(0x2F)
	prefixLength := 1

	if chunk[0] != prefix {
		return false
	}

	chunkLength := len(chunk)

	if chunkLength < 5 {
		return false
	}

	const numberOfBlocksLength = 4
	numberOfBlocks := ReadInt32(chunk[prefixLength : prefixLength+numberOfBlocksLength])

	weightSectionLength := 2 * numberOfBlocks
	offsetSectionLength := 4 * numberOfBlocks
	headerLength := prefixLength + numberOfBlocksLength + weightSectionLength + offsetSectionLength

	offsetSectionOffset := headerLength - offsetSectionLength
	firstOffset := ReadInt32(chunk[offsetSectionOffset : offsetSectionOffset+4])

	firstCodeBlockOffset := offsetSectionOffset + firstOffset + 4

	if chunkLength <= firstCodeBlockOffset {
		return false
	}

	firstCodeBlock := chunk[firstCodeBlockOffset:]
	expectedLength, ok := getExpectedLength(firstCodeBlock, firstCodeBlockOffset)

	if ok {
		return requirePrefixAndLength(prefix, expectedLength)(chunk)
	} else {
		return false
	}
}

func getExpectedLength(firstCodeBlock []byte, firstCodeBlockOffset int) (int, bool) {
	nextChunk := firstCodeBlock
	distanceTravelled := 0
	longJumpParameter := 0
	for {
		token, gotOne := searchForToken(nextChunk)
		distanceTravelled += len(token.Chunk)
		nextChunk = firstCodeBlock[distanceTravelled:]
		if gotOne {
			if token.Type == LongJump {
				longJumpParameter = ReadInt32(token.Chunk[1:])
				break
			}
		} else {
			return 0, false
		}
	}

	expectedLength := firstCodeBlockOffset + distanceTravelled + longJumpParameter
	return expectedLength, true
}
