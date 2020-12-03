package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/byxor/NeverScript/compiler"
	"os/exec"
	"path/filepath"
)

const (
	banner = `
     __                    __           _       _   
  /\ \ \_____   _____ _ __/ _\ ___ _ __(_)_ __ | |_ 
 /  \/ / _ \ \ / / _ \ '__\ \ / __| '__| | '_ \| __|
/ /\  /  __/\ V /  __/ |  _\ \ (__| |  | | |_) | |_ 
\_\ \/ \___| \_/ \___|_|  \__/\___|_|  |_| .__/ \__|
                                         |_|        
           The QB programming language.
----------------------------------------------------
`
	version = "0.4"
)

func main() {
	arguments := ParseCommandLineArguments()
	RunNeverscript(arguments)
}

func RunNeverscript(arguments CommandLineArguments) {
	argumentsWereSupplied := false

	if *arguments.FileToCompile != "" {
		argumentsWereSupplied = true

		var outputFilename string
		if *arguments.OutputFileName != "" {
			outputFilename = *arguments.OutputFileName
		} else {
			outputFilename = NsToQb(*arguments.FileToCompile)
		}

		fmt.Printf("\nCompiling '%s' (may freeze)...\n", *arguments.FileToCompile)
		var lexer compiler.Lexer
		var parser compiler.Parser
		var bytecodeCompiler compiler.BytecodeCompiler
		compiler.Compile(*arguments.FileToCompile, outputFilename, &lexer, &parser, &bytecodeCompiler)
		fmt.Printf("  Created '%s'.\n\n", outputFilename)

		if *arguments.ShowHexDump {
			fmt.Printf("Hex dump:\n%s", outputFilename, hex.Dump(bytecodeCompiler.Bytes))
		}

		if *arguments.DecompileWithRoq {
			fmt.Println("Roq decompiler output (may freeze):")
			roqCmd := exec.Command(".\\roq.exe", "-d", outputFilename)
			decompiledCode, _ := roqCmd.Output()
			fmt.Println(string(decompiledCode))
		}

		fmt.Println("done.")
	}

	if !argumentsWereSupplied {
		fmt.Println(banner[1:])
		fmt.Printf("Release %s\n\n", version)
		flag.Usage()
	}
}

type CommandLineArguments struct {
	FileToCompile    *string
	OutputFileName   *string
	ShowHexDump      *bool
	DecompileWithRoq *bool
}

func ParseCommandLineArguments() CommandLineArguments {
	args := CommandLineArguments{
		FileToCompile:   flag.String("c", "", "Specify a file to compile (.ns)."),
		OutputFileName:  flag.String("o", "", "Specify the output file name."),
		ShowHexDump:      flag.Bool("showHexDump", false, "Display the compiled bytecode in hex format."),
		DecompileWithRoq: flag.Bool("decompileWithRoq", false, "Display output from roq decompiler (roq.exe must be in your PATH)."),
	}
	flag.Parse()
	return args
}

func QbToNs(fileName string) string {
	return withoutExtension(fileName) + ".ns"
}

func NsToQb(fileName string) string {
	return withoutExtension(fileName) + ".qb"
}

func withoutExtension(fileName string) string {
	fileExtension := filepath.Ext(fileName)
	end := len(fileName)-len(fileExtension)
	return fileName[:end]
}
