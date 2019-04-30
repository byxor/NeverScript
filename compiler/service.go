package compiler

type Service interface {
	Compile(sourceCode SourceCode) (ByteCode, error)
}

type service struct {}

func NewService() Service {
	return &service{}
}

type SourceCode string
type ByteCode []byte
