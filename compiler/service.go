// Package compiler provides the use-case of compiling
// NeverScript source code into QB ByteCode.
//
// Used by modders.
package compiler

import (
	"github.com/byxor/NeverScript"
)

type Service interface {
	Compile(sourceCode NeverScript.SourceCode) (NeverScript.ByteCode, error)
}

// service is the implementation of Service.
type service struct {}

func (this service) Compile(sourceCode NeverScript.SourceCode) (NeverScript.ByteCode, error) {
	return NeverScript.NewByteCode([]byte{}), nil
}

func NewService() Service {
	return &service{}
}
