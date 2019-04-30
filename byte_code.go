package NeverScript

import (
	"bytes"
	"errors"
)

type ByteCode struct {
	content []byte
	length  int
}

func NewByteCode(content []byte) ByteCode {
	return ByteCode{
		content: content,
		length:  len(content),
	}
}

func (this ByteCode) GetSlice(startIndex, endIndex int) (ByteCode, error) {
	if startIndex < 0 {
		return nilByteCode, indexOutOfRange
	}

	if startIndex >= endIndex {
		return nilByteCode, indexOutOfRange
	}

	if endIndex > this.length {
		return nilByteCode, indexOutOfRange
	}

	content := this.content[startIndex:endIndex]
	return NewByteCode(content), nil
}

func (this ByteCode) IsEqualTo(other ByteCode) bool {
	return bytes.Equal(this.content, other.content)
}

func (this ByteCode) IsLongerThan(other ByteCode) bool {
	return this.length > other.length
}

func (this ByteCode) IsShorterThan(other ByteCode) bool {
	return this.length < other.length
}

func (this ByteCode) IsSameLengthAs(other ByteCode) bool {
	return this.length == other.length
}

func (this ByteCode) Contains(other ByteCode) bool {
	if other.IsLongerThan(this) {
		return false
	}

	iterateUpTo := this.length - other.length + 1

	for i := 0; i < iterateUpTo; i++ {
		sliceOfThis, _ := this.GetSlice(i, i + other.length)

		if sliceOfThis.IsEqualTo(other) {
			return true
		}
	}

	return false
}

var (
	indexOutOfRange = errors.New("Index is out of range")
	nilByteCode = NewByteCode([]byte{})
)
