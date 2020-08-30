package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/byxor/NeverScript/compiler"
	"os"
	"os/exec"
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
	version = "0.2"
)

func main() {
	//fileToCompile := "C:\\Users\\Brandon\\go\\src\\github.com\\byxor\\NeverScript\\docs\\neverscript-syntax.ns"
	//fileToDecompile := ""
	//outputFileName := ""
	//showHexDump := true
	//decompileWithRoq := true
	//arguments := commandLineArguments{
	//	FileToCompile:    &fileToCompile,
	//	FileToDecompile:  &fileToDecompile,
	//	OutputFileName:   &outputFileName,
	//	ShowHexDump:      &showHexDump,
	//	DecompileWithRoq: &decompileWithRoq,
	//}
	arguments := parseCommandLineArguments()
	argumentsWereSupplied := false

	if *arguments.FileToCompile != "" {
		argumentsWereSupplied = true

		inputFileName := *arguments.FileToCompile

		var outputFilename string
		if *arguments.OutputFileName != "" {
			outputFilename = *arguments.OutputFileName
		} else {
			outputFilename = NsToQb(inputFileName)
		}

		fmt.Printf("\nCompiling '%s' (may freeze)...\n", inputFileName)
		var lexer compiler.Lexer
		var newParser compiler.NewParser
		var bytecodeCompiler compiler.BytecodeCompiler
		compiler.Compile(inputFileName, outputFilename, &lexer, &newParser, &bytecodeCompiler)
		fmt.Printf("  Created '%s'.\n", outputFilename)

		if *arguments.ShowHexDump {
			fmt.Printf("\nHex dump:\n%s", outputFilename, hex.Dump(bytecodeCompiler.Bytes))
		}

		if *arguments.DecompileWithRoq {
			fmt.Println("\nRoq decompiler output (may freeze):")
			roqCmd := exec.Command(".\\roq.exe", "-d", outputFilename)
			decompiledCode, _ := roqCmd.Output()
			fmt.Println(string(decompiledCode))
		}

		fmt.Println()
		fmt.Println("done.")
	}

	if *arguments.FileToDecompile != "" {
		argumentsWereSupplied = true

		fmt.Printf("Decompiling '%s'...\n", *arguments.FileToDecompile)
		fmt.Println("This is not implemented yet, sorry.")

		fmt.Println()
		fmt.Println("done.")
	}

	if !argumentsWereSupplied {
		fmt.Println(banner[1:])
		fmt.Printf("Release %s", version)
		flag.Usage()
	}
}

type commandLineArguments struct {
	FileToCompile    *string
	FileToDecompile  *string
	OutputFileName   *string
	ShowHexDump      *bool
	DecompileWithRoq *bool
}

func parseCommandLineArguments() commandLineArguments {
	args := commandLineArguments{
		FileToCompile:   flag.String("c", "", "Specify a file to compile (.ns)."),
		FileToDecompile: flag.String("d", "", "Specify a file to decompile (.qb)."),
		OutputFileName:  flag.String("o", "", "Specify the output file name."),
		ShowHexDump:      flag.Bool("showHexDump", false, "Display the compiled bytecode in hex format"),
		DecompileWithRoq: flag.Bool("decompileWithRoq", false, "Display output from roq decompiler (roq.exe must be in your PATH)"),
	}
	flag.Parse()
	return args
}

func check(err error) {
	if err != nil {
		fmt.Println()
		fmt.Println()
		fmt.Println("SOMETHING WENT WRONG:")
		fmt.Println(err)
		os.Exit(1)
	}
}

func QbToNs(filename string) string {
	return withoutExtension(filename) + ".ns"
}

func NsToQb(filename string) string {
	return withoutExtension(filename) + ".qb"
}

func withoutExtension(filename string) string {
	return filename[:len(filename)-3]
}
