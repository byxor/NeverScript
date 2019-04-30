// Package compiler provides the use-case of compiling NeverScript source code into QB ByteCode.
//
// Used by modders.
package compiler

import (
	"github.com/byxor/NeverScript"
	"github.com/byxor/NeverScript/checksums"
	"github.com/pkg/errors"
	"strconv"
	"encoding/binary"
	goErrors "errors"
)

type Service interface {
	Compile(sourceCode NeverScript.SourceCode) (NeverScript.ByteCode, error)
}

// service is the implementation of Service.
type service struct {
	checksumService checksums.Service
}

func NewService(checksumService checksums.Service) Service {
	return &service{
		checksumService: checksumService,
	}
}

var junkByte = byte(0xF4)

func (this service) Compile(sourceCode NeverScript.SourceCode) (NeverScript.ByteCode, error) {
	byteCode := NeverScript.NewEmptyByteCode()

	syntaxTree, err := buildSyntaxTreeFrom(sourceCode)
	if err != nil {
		return byteCode, errors.Wrap(err, "Failed to build syntax tree")
	}

	byteCode.Push(NeverScript.EndOfLineToken)

	for _, declaration := range syntaxTree.Declarations {

		if declaration.EndOfLine != "" {
			byteCode.Push(NeverScript.EndOfLineToken)
			continue
		}

		if declaration.BooleanAssignment != nil {
			assignment := declaration.BooleanAssignment
			nameBytes := []byte{junkByte, junkByte, junkByte, junkByte}
			value := convertBooleanTextToByte(assignment.Value)

			byteCode.Push(NeverScript.NameToken)
			byteCode.Push(nameBytes...)
			byteCode.Push(NeverScript.EqualsToken)
			// The QB format has no Boolean type.
			// Instead, we use Ints with a value of 0 or 1.
			byteCode.Push(NeverScript.IntToken, value, 0, 0, 0)
			continue
		}

		if declaration.IntegerAssignment != nil {
			assignment := declaration.IntegerAssignment

			nameBytes := []byte{junkByte, junkByte, junkByte, junkByte}

			value, err := convertIntegerNodeToUint32(*assignment.Value)
			if err != nil {
				return byteCode, errors.Wrap(err, "Failed to convert integer node to uint32")
			}
			valueBytes := convertUint32ToLittleEndian(value)

			byteCode.Push(NeverScript.NameToken)
			byteCode.Push(nameBytes...)
			byteCode.Push(NeverScript.EqualsToken)
			byteCode.Push(NeverScript.IntToken)
			byteCode.Push(valueBytes...)
			continue
		}

		if declaration.StringAssignment != nil {
			assignment := declaration.StringAssignment

			nameBytes := []byte{junkByte, junkByte, junkByte, junkByte}

			quotedString := assignment.Value
			unquotedString := unquote(quotedString)

			lengthBytes := convertUint32ToLittleEndian(uint32(len(unquotedString)))
			stringBytes := []byte(unquotedString)
			nullTerminator := byte(0x00)

			byteCode.Push(NeverScript.NameToken)
			byteCode.Push(nameBytes...)
			byteCode.Push(NeverScript.EqualsToken)
			byteCode.Push(NeverScript.StringToken)
			byteCode.Push(lengthBytes...)
			byteCode.Push(stringBytes...)
			byteCode.Push(nullTerminator)
		}
	}

	byteCode.Push(NeverScript.EndOfFileToken)

	return byteCode, nil
}

func convertBooleanTextToByte(text string) byte {
	if text == "true" {
		return 0x01
	} else {
		return 0x00
	}
}

func convertIntegerNodeToUint32(node integer) (uint32, error) {
	text, base, nodeIsEmpty := (func() (text string, base int, nodeIsEmpty bool) {
		nodeIsEmpty = false

		if node.Base16 != "" {
			base = 16
			text = node.Base16[2:]
			return
		}

		if node.Base10 != "" {
			base = 10
			text = node.Base10
			return
		}

		if node.Base8 != "" {
			base = 8
			text = node.Base8[2:]
			return
		}

		if node.Base2 != "" {
			base = 2
			text = node.Base2[2:]
			return
		}

		nodeIsEmpty = true
		return
	})()

	if nodeIsEmpty {
		return 0, goErrors.New("Integer node is empty")
	}

	temp, err := strconv.ParseUint(text, base, 32)
	if err != nil {
		return 0, err
	}

	return uint32(temp), nil
}

func unquote(string string) string {
	return string[1:len(string)-1]
}

func convertUint32ToLittleEndian(value uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, value)
	return bytes
}

