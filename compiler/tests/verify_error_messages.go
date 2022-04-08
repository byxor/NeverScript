package main

import (
	"encoding/hex"
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

	check(IncompleteAssignment())
	check(EOFWhileScanningStringLiteral())
}

func IncompleteAssignment() error {
	var code string
	var err error

	expectedMessage := "Incomplete assignment"

	code = "x = "
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 1))
	if err != nil { return err }

	code = `x = 1
y = `
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 2))
	if err != nil { return err }

	code = `script Foo {
	var_x = 10
	var_y = 
}`
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 3))
	if err != nil { return err }

	code = `script Foo {}
script Bar {
	x =
}`
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 3))
	if err != nil { return err }

	code = `script Foo {}
script Bar {
	Foo x=5
	Foo x=
}`
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 4))
	if err != nil { return err }

	code = `script Foo {}
script Bar {
	Foo x=5 \
        yyyyyyyyyyyyyyyyyyyy=
}`
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 4))
	if err != nil { return err }

	code = `my_struct = {
	x =
}`
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 2))
	if err != nil { return err }

	code = `my_struct = {
	x = 1
	y = 2
	z = 3
	a = {
		b =
	}
}`
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 6))
	if err != nil { return err }

	code = `script Foo {
	x = {
		x =
	}
}`
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 3))
	if err != nil { return err }

	return nil
}

func EOFWhileScanningStringLiteral() error {
	var code string
	var err error

	expectedMessage := "EOF while scanning string literal"

	code = `"`
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 1))
	if err != nil { return err }

	code = "\"\n\n\n\n"
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 1))
	if err != nil { return err }

	code = "\n\""
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 2))
	if err != nil { return err }

	code = `
my_string = "this is ok"
my_string = "and so is this"
my_string = "but this is not`
	err = compileAndCheckError(code, checkMessageAndLineNumber(expectedMessage, 4))
	if err != nil { return err }

	return nil
}

func checkMessageAndLineNumber(
	expectedMessage string,
	expectedLineNumber int,
	//expectedColumnNumber int,
	) func(err compiler.Error, qbOutput []byte) error {
	return func(err compiler.Error, qbOutput []byte) error {
		if err == nil {
			errorMessage := fmt.Sprintf("expecting error for '%s' but got nothing", expectedMessage)
			if len(qbOutput) > 0 {
				errorMessage += "\n" + hex.Dump(qbOutput)
			}
			return errors.New(errorMessage)
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

func compileAndCheckError(code string, compilationErrorChecker func(compilationError compiler.Error, qbOutput []byte) error) error {
	tempDir, err := ioutil.TempDir(os.TempDir(), "neverscript-temporary-testing-tempDir")
	if err != nil { return err }
	defer os.RemoveAll(tempDir)

	ioutil.WriteFile(tempDir+"/code.ns", []byte(code), 0644)

	var lexer compiler.Lexer
	var parser compiler.Parser
	var bytecodeCompiler compiler.BytecodeCompiler
	bytecodeCompiler.RemoveChecksums = true
	compilationError := compiler.Compile(tempDir+"/code.ns", tempDir+"/code.qb", &lexer, &parser, &bytecodeCompiler)

	qbOutput, _ := ioutil.ReadFile(tempDir+"/code.qb")

	return compilationErrorChecker(compilationError, qbOutput)
}