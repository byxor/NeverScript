package main

import (
	"fmt"
	"github.com/byxor/NeverScript/compiler"
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

/*
 * Rather than writing a gigantic suite of tests like I did the past ~4 times I've started this project,
 * I'm just going to write a short test that warns me when the compiler is producing different output than before.
 *
 * I suspect that the change-detection will be easier to work with than a bunch of ever-changing tests, even if
 * a bit of human effort is required when reviewing the output each time.
 */

func main() {
	tempDir, err := ioutil.TempDir(os.TempDir(), "neverscript-temporary-testing-tempDir")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	nsPath := tempDir + "/code.ns"
	qbPath := tempDir + "/code.qb"

	// Create neverscript file
	ioutil.WriteFile(nsPath, []byte(code), 0644)

	// Compile
	fmt.Println("Compiling ns...")
	var lexer compiler.Lexer
	var parser compiler.NewParser
	var bytecodeCompiler compiler.BytecodeCompiler
	compiler.Compile(nsPath, qbPath, &lexer, &parser, &bytecodeCompiler)
	fmt.Println()

	// Decompile
	fmt.Println("Decompiling roq...")
	roqCmd := exec.Command(".\\roq.exe", "-d", qbPath)
	decompiledRoq, _ := roqCmd.Output()
	fmt.Println()

	// Compare results (the stupid variable name makes sure the results are aligned in the debugger)
	mractual := strings.Replace(string(decompiledRoq), "\r", "", -1)
	expected := strings.Replace(expectedDecompiledRoq, "\r", "", -1)
	banner := "--------------------------------------------------------"
	if mractual != expected {
		fmt.Println(banner)
		fmt.Println("WARNING: COMPILER OUTPUT CHANGED")
		fmt.Println(banner)
		fmt.Println("Expected:")
		fmt.Println(banner)
		fmt.Println(expected)
		fmt.Println(banner)
		fmt.Println("But got:")
		fmt.Println(banner)
		fmt.Println(mractual)
		fmt.Println(banner)
		fmt.Println("Diff:")
		fmt.Println(banner)

		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(expected),
			B:        difflib.SplitLines(mractual),
			FromFile: "Expected",
			ToFile:   "Actual",
			Context:  3,
		}
		text, _ := difflib.GetUnifiedDiffString(diff)
		fmt.Println(text)
	} else {
		fmt.Println(mractual)
		fmt.Println(banner)
		fmt.Println("Hooray! Still working!")
	}
}

const code = `
my_int = 10
my_int = -10
my_float = 0.1
my_float = -0.1
my_string = "hey"
my_pair = (1.0, 2.0)
my_vector = (100.0, 200.0, 300.0)

script TestBasicExpressions {
    x = 1
	x = (1)
	
	description = "Positive ints"
	x = (1 + 2)
	x = (1 - 2)
    x = (1 * 3)
    x = (1 / 2)
    
	description = "Negative ints"
	x = (-1 + -2)
	x = (-1 - -2)
    x = (-1 * -3)
    x = (-1 / -2)
	
	description = "Positive floats"
	x = (1.0 + 2.0)
	x = (1.0 - 2.0)
    x = (1.0 * 3.0)
    x = (1.0 / 2.0)

	description = "Negative floats"
	x = (-1.0 + -2.0)
	x = (-1.0 - -2.0)
    x = (-1.0 * -3.0)
    x = (-1.0 / -2.0)
}

script TestShorthandMath {}
`

const expectedDecompiledRoq = `
:i $my_int$ = %i(10,0000000a)
:i $my_int$ = %i(4294967286,fffffff6)
:i $my_float$ = %f(0.100000)
:i $my_float$ = %f(-0.100000)
:i $my_string$ = %s(3,"hey")
:i $my_pair$ = %vec2(1.000000,2.000000)
:i $my_vector$ = %vec3(100.000000,200.000000,300.000000)
:i function $TestBasicExpressions$
	:i $x$ = %i(1,00000001)
	:i $x$ =  (%i(1,00000001)) 
	:i $description$ = %s(13,"Positive ints")
	:i $x$ =  (%i(1,00000001) + %i(2,00000002)) 
	:i $x$ =  (%i(1,00000001) - %i(2,00000002)) 
	:i $x$ =  (%i(1,00000001) * %i(3,00000003)) 
	:i $x$ =  (%i(1,00000001) / %i(2,00000002)) 
	:i $description$ = %s(13,"Negative ints")
	:i $x$ =  (%i(4294967295,ffffffff) + %i(4294967294,fffffffe)) 
	:i $x$ =  (%i(4294967295,ffffffff) - %i(4294967294,fffffffe)) 
	:i $x$ =  (%i(4294967295,ffffffff) * %i(4294967293,fffffffd)) 
	:i $x$ =  (%i(4294967295,ffffffff) / %i(4294967294,fffffffe)) 
	:i $description$ = %s(15,"Positive floats")
	:i $x$ =  (%f(1.000000) + %f(2.000000)) 
	:i $x$ =  (%f(1.000000) - %f(2.000000)) 
	:i $x$ =  (%f(1.000000) * %f(3.000000)) 
	:i $x$ =  (%f(1.000000) / %f(2.000000)) 
	:i $description$ = %s(15,"Negative floats")
	:i $x$ =  (%f(-1.000000) + %f(-2.000000)) 
	:i $x$ =  (%f(-1.000000) - %f(-2.000000)) 
	:i $x$ =  (%f(-1.000000) * %f(-3.000000)) 
	:i $x$ =  (%f(-1.000000) / %f(-2.000000)) 
:i endfunction
:i function $TestShorthandMath$endfunction
:i :end
`