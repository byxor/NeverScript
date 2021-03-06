package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/byxor/NeverScript/compiler"
	"github.com/byxor/NeverScript/decompiler"
	"github.com/byxor/NeverScript/pre_generator"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	usage = `
COMPILATION:
    -c                 (required string)  Specify a file to compile (.ns).
    -o                 (optional string)  Specify the output file name (.qb).
    -showHexDump       (optional flag)    Display the compiled bytecode in hex format.
    -decompileWithRoq  (optional flag)    Display output from roq decompiler (roq.exe must be in your PATH).

PRE GENERATION:
    -p                 (required string)  Specify a pre spec file (.ps).
    -showHexDump       (optional flag)    Display the pre bytes in hex format.

DECOMPILATION (very incomplete):
    -d                 (required string)  Specify a file to decompile (.qb).
    -o                 (optional string)  Specify the output file name (.ns).
    -showCode          (optional flag)    Display the decompiled code as text.
`

	version = "0.6"
)

type CommandLineArguments struct {
	FileToCompile    *string
	FileToDecompile  *string
	PreSpecFile      *string
	OutputFileName   *string
	ShowHexDump      *bool
	ShowCode         *bool
	DecompileWithRoq *bool
}

func main() {
	arguments := ParseCommandLineArguments()
	// Hardcoded arguments for testing:
	/* if len(os.Args) == 1 {
		*arguments.FileToCompile = "C:\\Users\\Brandon\\Desktop\\mod\\foo.ns" // build/PRE3,thugpro_qb.prx/qb/_mods/byxor_debug.qb"
		*arguments.ShowCode = true
		RunNeverscript(arguments)
		*arguments.FileToCompile = ""
		*arguments.FileToDecompile = "C:\\Users\\Brandon\\Desktop\\mod\\foo.qb" // build/PRE3,thugpro_qb.prx/qb/_mods/byxor_debug.qb"
		*arguments.ShowCode = true
		*arguments.OutputFileName = "nul"
	} */
	RunNeverscript(arguments)
}

func ParseCommandLineArguments() CommandLineArguments {
	args := CommandLineArguments{
		FileToCompile:    flag.String("c", "", ""),
		FileToDecompile:  flag.String("d", "", ""),
		PreSpecFile:      flag.String("p", "", ""),
		OutputFileName:   flag.String("o", "", ""),
		ShowHexDump:      flag.Bool("showHexDump", false, ""),
		ShowCode:      flag.Bool("showCode", false, ""),
		DecompileWithRoq: flag.Bool("decompileWithRoq", false, ""),
	}
	flag.Parse()
	return args
}

func RunNeverscript(arguments CommandLineArguments) {
	argumentsWereSupplied := false

	if *arguments.FileToCompile != "" {
		argumentsWereSupplied = true

		outputFilename := *arguments.OutputFileName
		if outputFilename == "" {
			outputFilename = WithQbExtension(*arguments.FileToCompile)
		}

		fmt.Printf("\nCompiling '%s' (may freeze)...\n", *arguments.FileToCompile)
		var lexer compiler.Lexer
		var parser compiler.Parser
		var bytecodeCompiler compiler.BytecodeCompiler
		compiler.Compile(*arguments.FileToCompile, outputFilename, &lexer, &parser, &bytecodeCompiler)
		fmt.Printf("  Created '%s'.\n", outputFilename)

		if *arguments.ShowHexDump {
			fmt.Printf("\n%s\n", hex.Dump(bytecodeCompiler.Bytes))
		} else {
			fmt.Println()
		}

		if *arguments.DecompileWithRoq {
			fmt.Println("Roq decompiler output (may freeze):")
			roqCmd := exec.Command("roq.exe", "-d", outputFilename)
			decompiledCode, _ := roqCmd.Output()
			fmt.Println(string(decompiledCode))
		}
	} else if *arguments.FileToDecompile != "" {
		argumentsWereSupplied = true

		fmt.Printf("\nDecompiling '%s' (may freeze)...\n", *arguments.FileToDecompile)
		byteCode, err := ioutil.ReadFile(*arguments.FileToDecompile)
		if err != nil {
			log.Fatal(err)
		}

		var decompilerArguments decompiler.Arguments
		decompilerArguments.ByteCode = byteCode
		decompiler.Decompile(&decompilerArguments)

		outputFilename := *arguments.OutputFileName
		if outputFilename == "" {
			outputFilename = WithNsExtension(*arguments.FileToDecompile)
		}

		ioutil.WriteFile(outputFilename, []byte(decompilerArguments.SourceCode), 0644)
		fmt.Printf("    Created '%s'.\n", outputFilename)

		if *arguments.ShowCode {
			fmt.Printf("\n```%s\n```\n", decompilerArguments.SourceCode)
		} else {
			fmt.Println()
		}
	} else if *arguments.PreSpecFile != "" {
		argumentsWereSupplied = true

		outputFilename := *arguments.OutputFileName
		if outputFilename == "" {
			outputFilename = WithPrxExtension(*arguments.PreSpecFile)
		}

		fmt.Printf("\nGenerating pre file from spec '%s'...\n", *arguments.PreSpecFile)
		preSpec := pre_generator.ParsePreSpec(*arguments.PreSpecFile)
		pre := pre_generator.MakePre(preSpec)
		ioutil.WriteFile(*arguments.PreSpecFile, pre, 0466)
		fmt.Printf("  Created '%s'.\n\n", outputFilename)

		if *arguments.ShowHexDump {
			fmt.Printf("Hex dump:\n%s\n", hex.Dump(pre))
		}
	}

	if !argumentsWereSupplied {
		fmt.Println(banner[1:])
		fmt.Printf("Release %s\n\n", version)
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf(usage)
	} else {
		banner := strings.Repeat("-", 36)
		fmt.Println(banner + " done " + banner)
	}
}

func WithQbExtension(fileName string) string {
	return withoutExtension(fileName) + ".qb"
}

func WithPrxExtension(fileName string) string {
	return withoutExtension(fileName) + ".prx"
}

func WithNsExtension(fileName string) string {
	return withoutExtension(fileName) + ".ns"
}

func withoutExtension(fileName string) string {
	fileExtension := filepath.Ext(fileName)
	end := len(fileName) - len(fileExtension)
	return fileName[:end]
}
