package NeverScript

import (
	"bytes"
	goErrors "errors"
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

func NewEmptyByteCode() ByteCode {
	return NewByteCode([]byte{})
}

func (this *ByteCode) Clear() {
	this.content = []byte{}
	this.length = 0
}

func (this *ByteCode) Push(bytes ...byte) {
	this.content = append(this.content, bytes...)
	this.length += len(bytes)
}

func (this ByteCode) GetSlice(startIndex, endIndex int) (ByteCode, error) {
	if startIndex < 0 {
		return nilByteCode, SliceIndexOutOfRange
	}

	if startIndex >= endIndex {
		return nilByteCode, SliceIndexOutOfRange
	}

	if endIndex > this.length {
		return nilByteCode, SliceIndexOutOfRange
	}

	content := this.content[startIndex:endIndex]
	return NewByteCode(content), nil
}

func (this ByteCode) ToBytes() []byte {
	return this.content
}

func (this ByteCode) IsLongerThan(other ByteCode) bool {
	return this.length > other.length
}

func (this ByteCode) IsEqualTo(other ByteCode) bool {
	return bytes.Equal(this.content, other.content)
}

func (this ByteCode) Contains(other ByteCode) (bool, error) {
	if other.IsLongerThan(this) {
		return false, nil
	}

	iterateUpTo := this.length - other.length + 1

	for i := 0; i < iterateUpTo; i++ {
		// Error can be ignored as it's impossible to go out of range
		sliceOfThis, _ := this.GetSlice(i, i+other.length)

		if sliceOfThis.IsEqualTo(other) {
			return true, nil
		}
	}

	return false, nil
}

func (this ByteCode) Contains_IgnoreByte(other ByteCode, byteToIgnore byte) (bool, error) {
	forceFutureComparisonsToPass(&this, other, byteToIgnore)
	return this.Contains(other)
}

func forceFutureComparisonsToPass(byteCodeToModify *ByteCode, byteCodeToCheck ByteCode, byteToIgnore byte) {
	length := min(byteCodeToCheck.length, byteCodeToModify.length)

	for i := 0; i < length; i++ {
		byteToCheck := byteCodeToCheck.content[i]

		if byteToCheck == byteToIgnore {
			// This will force the comparison checks to pass.
			byteCodeToModify.content[i] = byteToIgnore
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

var (
	SliceIndexOutOfRange = goErrors.New("Index is out of range")
	nilByteCode          = NewEmptyByteCode()
)
