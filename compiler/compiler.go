package compiler

import (
	"io/ioutil"
	"log"
	"strings"
)

func Compile(nsFilePath, qbFilePath string, lexer *Lexer, newParser *NewParser, bytecodeCompiler *BytecodeCompiler) {
	{ // read source code into memory & store it in lexer
		bytes, err := ioutil.ReadFile(nsFilePath)
		if err != nil {
			log.Fatal(err)
		}
		lexer.SourceCode = string(bytes)

		// Remove weird windows line-endings
		lexer.SourceCode = strings.Replace(lexer.SourceCode, "\r", "", -1)

		lexer.SourceCodeSize = len(lexer.SourceCode)
	}

	LexSourceCode(lexer)

	newParser.Tokens = lexer.Tokens
	BuildAbstractSyntaxTree(newParser)
	if !newParser.Result.WasSuccessful {
		log.Fatal(newParser.Result.Reason)
	}

	bytecodeCompiler.RootAstNode = newParser.Result.Node
	GenerateBytecode(bytecodeCompiler)

	ioutil.WriteFile(qbFilePath, bytecodeCompiler.Bytes, 0644)
}
