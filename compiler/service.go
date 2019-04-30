// Package compiler provides the use-case of compiling
// NeverScript source code into QB ByteCode.
//
// Used by modders.
package compiler

import (
	"github.com/byxor/NeverScript"
	"github.com/byxor/NeverScript/checksums"
	"github.com/pkg/errors"
	"log"
	"strconv"
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
		log.Printf("%+v\n", declaration)

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

			checksum, err := convertIntegerNodeToChecksum(*assignment.Value)

			if err != nil {
				return byteCode, errors.Wrap(err, "Failed to convert compiler.integer to checksum")
			}

			valueBytes := this.checksumService.EncodeAsLittleEndian(checksum)

			byteCode.Push(NeverScript.NameToken)
			byteCode.Push(nameBytes...)
			byteCode.Push(NeverScript.EqualsToken)
			byteCode.Push(NeverScript.IntToken)
			byteCode.Push(valueBytes...)
			continue
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

func convertIntegerNodeToChecksum(node integer) (NeverScript.Checksum, error) {
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
		return NeverScript.NewEmptyChecksum(), goErrors.New("Integer node is empty")
	}

	temp, err := strconv.ParseUint(text, base, 32)
	if err != nil {
		return NeverScript.NewEmptyChecksum(), err
	}

	content := uint32(temp)

	return NeverScript.NewChecksum(content), nil
}
