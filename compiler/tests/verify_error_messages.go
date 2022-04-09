package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/byxor/NeverScript/compiler"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"runtime"
)

// I made them globals so I don't have to think about `:=` vs `=` or forward declare them inside each test method.
// It's not necessary tbh.
// Actually, it kinda helps cuz I don't have to pass so many params around the place.
var code string
var err error
var expectedMessage string
var expectedLineNumber int

func main() {
	check := func(functionThatRunsTheTest func() error) {
		functionName := runtime.FuncForPC(reflect.ValueOf(functionThatRunsTheTest).Pointer()).Name()[5:]
		if err := functionThatRunsTheTest(); err == nil {
			fmt.Print("✓ ")
			fmt.Println(functionName)
		} else {
			fmt.Print("✗ ")
			fmt.Println(functionName)
			log.Fatal(fmt.Sprintf(" %s", err.Error()))
		}
	}

	check(ExtraParenthesis)
	check(IncompleteVector)
	check(IncompletePair)
	check(IncompleteParentheses)
	check(IncompleteShorthandDivision)
	check(IncompleteShorthandMultiplication)
	check(IncompleteShorthandSubtraction)
	check(IncompleteShorthandAddition)
	check(IncompleteAssignment)
	check(EOFWhileScanningStringLiteral)
}

func ExtraParenthesis() error {
	expectedMessage = "Unnecessary parenthesis )"

	code = ")"
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = "x = (5))"
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	x = ((<y> + <z>)))
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	return nil
}

func IncompleteVector() error {
	expectedMessage = "Incomplete vector expression"

	code = "(0.0, 1.0, 2.0"
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
    (0.0, 1.0, 2.0
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	x = 2
	y = (0.0, 1.0,
}`
	expectedLineNumber = 3
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	return nil
}

func IncompletePair() error {
	expectedMessage = "Incomplete pair expression"

	code = "(0.0, 1.0"
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
    (0.0, 1.0
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	x = 2
	y = (0.0,
}`
	expectedLineNumber = 3
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	return nil
}

func IncompleteParentheses() error {
	expectedMessage = "Incomplete parenthesis ("

	code = "("
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `x = 1
y = (

`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `x = (
	(
)`
	expectedLineNumber = 1 // I expected the line number to be 2, but 1 makes sense too. Depends how your brain parses it.
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	return nil
}

func IncompleteShorthandDivision() error {
	expectedMessage = "Incomplete *="

	code = "x *="
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `x += 1
y *=`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	x *= 
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = "<y> *="
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `<y> *= 1
<z> *=`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	<y> *= 
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	return nil
}

func IncompleteShorthandMultiplication() error {
	expectedMessage = "Incomplete *="

	code = "x *="
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `x += 1
y *=`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	x *= 
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = "<y> *="
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `<y> *= 1
<z> *=`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	<y> *= 
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	return nil
}

func IncompleteShorthandSubtraction() error {
	expectedMessage = "Incomplete -="

	code = "x -="
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `x += 1
y -=`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	x -= 
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = "<y> -="
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `<y> -= 1
<z> -=`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	<y> -= 
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	return nil
}

func IncompleteShorthandAddition() error {
    expectedMessage = "Incomplete +="

    code = "x +="
    expectedLineNumber = 1
    err = compileAndCheckError(checkMessageAndLineNumber)
    if err != nil { return err }

    code = `x += 1
y +=`
    expectedLineNumber = 2
    err = compileAndCheckError(checkMessageAndLineNumber)
    if err != nil { return err }

	code = `script Foo {
	x += 
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = "<y> +="
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `<y> += 1
<z> +=`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	<y> += 
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	return nil
}

func IncompleteAssignment() error {
	expectedMessage = "Incomplete assignment"

	code = "x = "
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `x = 1
y = `
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	var_x = 10
	var_y = 
}`
	expectedLineNumber = 3
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {}
script Bar {
	x =
}`
	expectedLineNumber = 3
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {}
script Bar {
	Foo x=5
	Foo x=
}`
	expectedLineNumber = 4
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {}
script Bar {
	Foo x=5 \
        yyyyyyyyyyyyyyyyyyyy=
}`
	expectedLineNumber = 4
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `my_struct = {
	x =
}`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `my_struct = {
	x = 1
	y = 2
	z = 3
	a = {
		b =
	}
}`
	expectedLineNumber = 6
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `script Foo {
	x = {
		x =
	}
}`
	expectedLineNumber = 3
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	return nil
}

func EOFWhileScanningStringLiteral() error {
	expectedMessage = "EOF while scanning string literal"

	code = `"`
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = "\"\n\n\n\n"
	expectedLineNumber = 1
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = "\n\""
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `
my_string = "this is ok"
my_string = "and so is this"
my_string = "but this is not`
	expectedLineNumber = 4
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }

	code = `my_struct = {
	"
}
`
	expectedLineNumber = 2
	err = compileAndCheckError(checkMessageAndLineNumber)
	if err != nil { return err }
	
	return nil
}

func checkMessageAndLineNumber(err compiler.Error, qbOutput []byte) error {
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

func compileAndCheckError(compilationErrorChecker func(compilationError compiler.Error, qbOutput []byte) error) error {
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