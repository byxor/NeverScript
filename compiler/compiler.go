package compiler

import (
	"io/ioutil"
	"log"
	"strings"
)

func Compile(nsFilePath, qbFilePath string, lexer *Lexer, parser *Parser, bytecodeCompiler *BytecodeCompiler) Error {
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

	err := LexSourceCode(lexer)
	if err != nil { return err }

	parser.Tokens = lexer.Tokens
	BuildAbstractSyntaxTree(parser)
	if !parser.Result.GotResult {
		log.Fatal(parser.Result.Reason)
	} else if parser.Result.Error != nil {
		return CompilationError{
			message:      parser.Result.Error.Error(),
			lineNumber:   parser.Result.LineNumber,
		}
	}

	bytecodeCompiler.RootAstNode = parser.Result.Node
	GenerateBytecode(bytecodeCompiler)

	ioutil.WriteFile(qbFilePath, bytecodeCompiler.Bytes, 0644)

	return nil
}
