package main

import (
	"errors"
	"fmt"
	"github.com/byxor/NeverScript/compiler"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	check := func(err error) {
		if err == nil {
			fmt.Print("✓")
		} else {
			log.Fatal(err)
		}
	}

	check(EOFWhileScanningStringLiteral())
}

func EOFWhileScanningStringLiteral() error {
	var code string
	var err error

	expectedMessage := "EOF while scanning string literal"

	code = `"`
	err = compileAndCheckError(code, check(expectedMessage, 1))
	if err != nil { return err }

	code = "\"\n\n\n\n"
	err = compileAndCheckError(code, check(expectedMessage, 1))
	if err != nil { return err }

	code = "\n\""
	err = compileAndCheckError(code, check(expectedMessage, 2))
	if err != nil { return err }

	code = `
my_string = "this is ok"
my_string = "and so is this"
my_string = "but this is not`
	err = compileAndCheckError(code, check(expectedMessage, 4))
	if err != nil { return err }

	return nil
}

func check(
	expectedMessage string,
	expectedLineNumber int,
	//expectedColumnNumber int,
	) func(err compiler.Error) error {
	return func(err compiler.Error) error {
		if err == nil {
			return errors.New("expecting error for incomplete string")
		}
		if message := err.GetMessage(); message != expectedMessage {
			return errors.New(fmt.Sprintf("expecting message to be '%s' but got '%s'", expectedMessage, message))
		}
		if lineNumber := err.GetLineNumber(); lineNumber != expectedLineNumber {
			return errors.New(fmt.Sprintf("expecting line number to be %d but got %d", expectedLineNumber, lineNumber))
		}
		//if columnNumber := err.GetColumnNumber(); columnNumber != expectedColumnNumber {
		//	return errors.New(fmt.Sprintf("expecting column number to be %d but got %d", expectedColumnNumber, columnNumber))
		//}
		return nil
	}
}

func compileAndCheckError(code string, compilationErrorChecker func(compilationError compiler.Error) error) error {
	tempDir, err := ioutil.TempDir(os.TempDir(), "neverscript-temporary-testing-tempDir")
	if err != nil { return err }
	defer os.RemoveAll(tempDir)

	ioutil.WriteFile(tempDir+"/code.ns", []byte(code), 0644)

	var lexer compiler.Lexer
	var parser compiler.Parser
	var bytecodeCompiler compiler.BytecodeCompiler
	compilationError := compiler.Compile(tempDir+"/code.ns", tempDir+"/code.qb", &lexer, &parser, &bytecodeCompiler)

	return compilationErrorChecker(compilationError)
}