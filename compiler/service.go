// Package compiler provides the use-case of compiling NeverScript source code into QB ByteCode.
//
// Used by modders.
package compiler

import (
	"encoding/binary"
	"fmt"
	goErrors "errors"
	"github.com/byxor/NeverScript"
	"github.com/byxor/NeverScript/checksums"
	"github.com/pkg/errors"
	"strconv"
)

type Service interface {
	Compile(sourceCode NeverScript.SourceCode) (NeverScript.ByteCode, error)
}

// service is the implementation of Service.
type service struct {
	checksumService checksums.Service
	byteCode        NeverScript.ByteCode
}

func NewService(checksumService checksums.Service) Service {
	return &service{
		checksumService: checksumService,
		byteCode:        NeverScript.NewEmptyByteCode(),
	}
}

var junkByte = byte(0xF4)

func (this *service) Compile(sourceCode NeverScript.SourceCode) (NeverScript.ByteCode, error) {
	this.byteCode.Clear()

	syntaxTree, err := buildSyntaxTreeFrom(sourceCode)
	if err != nil {
		return this.byteCode, errors.Wrap(err, "Failed to build syntax tree")
	}

	this.byteCode.Push(NeverScript.EndOfLineToken)

	for _, declaration := range syntaxTree.Declarations {
		this.processDeclaration(*declaration)
	}

	this.byteCode.Push(NeverScript.EndOfFileToken)

	return this.byteCode, nil
}

func (this *service) processDeclaration(declaration declaration) error {
	if declaration.EndOfLine != "" {
		this.byteCode.Push(NeverScript.EndOfLineToken)
		return nil
	}

	if declaration.BooleanAssignment != nil {
		err := this.processBooleanAssignment(*declaration.BooleanAssignment)
		return errors.Wrap(err, "Failed to process Boolean Assignment")
	}

	if declaration.IntegerAssignment != nil {
		err := this.processIntegerAssignment(*declaration.IntegerAssignment)
		return errors.Wrap(err, "Failed to process Integer Assignment")
	}

	if declaration.StringAssignment != nil {
		this.processStringAssignment(*declaration.StringAssignment)
		return nil
	}

	return nil
}

func (this *service) processBooleanAssignment(assignment booleanAssignment) error {
	nameBytes := []byte{junkByte, junkByte, junkByte, junkByte}

	value, err := convertBooleanTextToByte(assignment.Value)
	if err !=nil {
		return errors.Wrap(err, "Failed to convert boolean text to byte")
	}

	this.byteCode.Push(NeverScript.NameToken)
	this.byteCode.Push(nameBytes...)
	this.byteCode.Push(NeverScript.EqualsToken)
	// The QB format has no Boolean type.
	// Instead, we use Ints with a value of 0 or 1.
	this.byteCode.Push(NeverScript.IntToken, value, 0, 0, 0)

	return nil
}

func (this *service) processIntegerAssignment(assignment integerAssignment) error {
	nameBytes := []byte{junkByte, junkByte, junkByte, junkByte}

	value, err := convertIntegerNodeToUint32(*assignment.Value)
	if err != nil {
		return errors.Wrap(err, "Failed to convert integer node to uint32")
	}
	valueBytes := convertUint32ToLittleEndian(value)

	this.byteCode.Push(NeverScript.NameToken)
	this.byteCode.Push(nameBytes...)
	this.byteCode.Push(NeverScript.EqualsToken)
	this.byteCode.Push(NeverScript.IntToken)
	this.byteCode.Push(valueBytes...)

	return nil
}

func (this *service) processStringAssignment(assignment stringAssignment) {
	nameBytes := []byte{junkByte, junkByte, junkByte, junkByte}

	quotedString := assignment.Value
	unquotedString := unquote(quotedString)

	lengthBytes := convertUint32ToLittleEndian(uint32(len(unquotedString)))
	stringBytes := []byte(unquotedString)
	nullTerminator := byte(0x00)

	this.byteCode.Push(NeverScript.NameToken)
	this.byteCode.Push(nameBytes...)
	this.byteCode.Push(NeverScript.EqualsToken)
	this.byteCode.Push(NeverScript.StringToken)
	this.byteCode.Push(lengthBytes...)
	this.byteCode.Push(stringBytes...)
	this.byteCode.Push(nullTerminator)
}

func convertBooleanTextToByte(text string) (byte, error) {
	if text == "true" {
		return 0x01, nil
	} else if text == "false" {
		return 0x00, nil
	}
	return 0, goErrors.New(fmt.Sprintf("Cannot convert '%s' to a 0 or 1.", text))
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
	return string[1 : len(string)-1]
}

func convertUint32ToLittleEndian(value uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, value)
	return bytes
}
