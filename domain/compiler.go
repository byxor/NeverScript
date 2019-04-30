package domain

type ByteCode []byte

type Compiler interface {
	GenerateByteCode(code string) (ByteCode, error)
}
