package compiler

type Error interface {
	GetMessage() string
	GetLineNumber() int
	GetColumnNumber() int
}

type CompilationError struct {
	message string
	lineNumber int
	columnNumber int
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