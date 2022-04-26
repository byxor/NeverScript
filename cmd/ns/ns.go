package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"github.com/byxor/NeverScript/compiler"
	"github.com/byxor/NeverScript/decompiler"
	"github.com/byxor/NeverScript/pre_generator"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
    -targetGame        (optional string)  Specify which game to target (defaults to "thug2").
    -removeChecksums   (optional flag)    Removes checksum information from end of output.
    -showHexDump       (optional flag)    Display the compiled bytecode in hex format.
    -showDecompiledRoq (optional flag)    Display output from roq decompiler (roq.exe must be in your PATH).

PRE GENERATION:
    -p                 (required string)  Specify a pre spec file (.ps).
    -showHexDump       (optional flag)    Display the pre bytes in hex format.

DECOMPILATION (very incomplete):
    -d                 (required string)  Specify a file to decompile (.qb).
    -o                 (optional string)  Specify the output file name (.ns).
    -showCode          (optional flag)    Display the decompiled code as text.

`

	version = "0.9"
)

type CommandLineArguments struct {
	FileToCompile     *string
	FileToDecompile   *string
	PreSpecFile       *string
	OutputFileName    *string
	TargetGame        *string
	ShowHexDump       *bool
	ShowCode          *bool
	RemoveChecksums   *bool
	ShowDecompiledRoq *bool
}

func main() {
	arguments := ParseCommandLineArguments()
	if err := RunNeverscript(arguments); err != nil {
		fmt.Println()
	    fmt.Println(err.Error())
		os.Exit(1)
	}
}

func ParseCommandLineArguments() CommandLineArguments {
	args := CommandLineArguments{
		FileToCompile:     flag.String("c", "", ""),
		FileToDecompile:   flag.String("d", "", ""),
		PreSpecFile:       flag.String("p", "", ""),
		OutputFileName:    flag.String("o", "", ""),
		TargetGame:        flag.String("targetGame", "thug2", ""),
		ShowHexDump:       flag.Bool("showHexDump", false, ""),
		ShowCode:          flag.Bool("showCode", false, ""),
		ShowDecompiledRoq: flag.Bool("showDecompiledRoq", false, ""),
		RemoveChecksums:   flag.Bool("removeChecksums", false, ""),
	}
	flag.Parse()
	return args
}

func RunNeverscript(arguments CommandLineArguments) error {
	argumentsWereSupplied := false

	//fileToDecompile := `C:\Program Files (x86)\Aspyr\Tony Hawks Pro Skater 4\Game\data\scripts\Levels.qb`
	//showCode := true
	//arguments.FileToDecompile = &fileToDecompile
	//arguments.ShowCode = &showCode

	if *arguments.FileToCompile != "" {
		argumentsWereSupplied = true

		outputFileName := *arguments.OutputFileName
		if outputFileName == "" {
			outputFileName = WithQbExtension(*arguments.FileToCompile)
		}

		var lexer compiler.Lexer
		var parser compiler.Parser
		var bytecodeCompiler compiler.BytecodeCompiler
		bytecodeCompiler.TargetGame = strings.ToLower(*arguments.TargetGame)
		bytecodeCompiler.RemoveChecksums = *arguments.RemoveChecksums

		if bytecodeCompiler.TargetGame != "thps3" &&
			bytecodeCompiler.TargetGame != "thps4" &&
			bytecodeCompiler.TargetGame != "thug1" &&
			bytecodeCompiler.TargetGame != "thug2" {
			return errors.New("ERROR - Target game must be thps3/thps4/thug1/thug2")
		}

		compilationChannel := make(chan compiler.Error, 1)
		go func() {
			compilationError := compiler.Compile(*arguments.FileToCompile, outputFileName, &lexer, &parser, &bytecodeCompiler)
			compilationChannel <- compilationError
		}()
		var compilationError compiler.Error
		select {
		case result := <-compilationChannel:
		    compilationError = result
		case <-time.After(3 * time.Second):
			return errors.New("ERROR - Compiler took too long. It probably went into an infinite loop because of a bug or an unimplemented feature")
			fmt.Println("\nWARNING - Roq decompiler froze. Some QB cannot be decompiled, e.g. adjacent line ending bytes (0x01 0x01)")
		}
		if compilationError != nil {
			return compilationError.ToError()
		}
		fmt.Printf("\n  Created '%s'.\n", outputFileName)

		if *arguments.ShowHexDump {
			fmt.Printf("\n%s", hex.Dump(bytecodeCompiler.Bytes))
		}

		if *arguments.ShowDecompiledRoq {
			fmt.Println("\nRoq decompiler output:")

			roqChannel := make(chan string, 1)

			go func() {
				roqCmd := exec.Command("roq.exe", "-d", outputFileName)
				decompiledCode, _ := roqCmd.Output()
				roqChannel <- "\n" + strings.TrimSpace(string(decompiledCode))
			}()

			select {
			case decompiledRoq := <-roqChannel:
				fmt.Println(decompiledRoq)
			case <-time.After(3 * time.Second):
				fmt.Println("\nWARNING - Roq decompiler froze. Some QB cannot be decompiled, e.g. adjacent line ending bytes (0x01 0x01)")
			}
		}
	} else if *arguments.FileToDecompile != "" {
		argumentsWereSupplied = true

		qb, err := ioutil.ReadFile(*arguments.FileToDecompile)
		if err != nil {
			return err
		}

		decompiledCode, err := decompiler.Decompile(qb)
		if err != nil {
			return err
		}

		outputFileName := *arguments.OutputFileName
		if outputFileName == "" {
			outputFileName = WithNsExtension(*arguments.FileToDecompile)
		}
		ioutil.WriteFile(outputFileName, []byte(decompiledCode), 0644)

		fmt.Printf("\n  Created '%s'.\n", outputFileName)

		if *arguments.ShowCode {
			fmt.Printf("\n%s", decompiledCode)
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
		ioutil.WriteFile(*arguments.PreSpecFile, pre, 0644)
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
	}

	return nil
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
