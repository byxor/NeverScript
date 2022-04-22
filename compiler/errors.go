package compiler

import (
	"errors"
	"fmt"
)

type Error interface {
	GetMessage() string
	GetLineNumber() int
	GetColumnNumber() int
	ToError() error
}

type CompilationError struct {
	message string
	lineNumber int
	columnNumber int
	baseFilePath string
}

func (self CompilationError) GetMessage() string {
	return self.message
}

func (self CompilationError) GetLineNumber() int {
	return self.lineNumber
}

func (self CompilationError) GetColumnNumber() int {
	return self.columnNumber
}

func (self CompilationError) ToError() error {
	return errors.New(fmt.Sprintf("ERROR %s(line %d) - %s", self.baseFilePath, self.lineNumber, self.message))
}