package cmd

import (
	"flag"
	"fmt"
	"github.com/byxor/NeverScript"
	"github.com/byxor/NeverScript/compiler"
	"github.com/byxor/NeverScript/shared/filenames"
	"io/ioutil"
	"os"
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
)

var compilerService = compiler.NewService()

func main() {
	arguments := parseCommandLineArguments()
	argumentsWereSupplied := false

	fmt.Println(banner[1:])

	if *arguments.FileToCompile != "" {
		argumentsWereSupplied = true

		inputFileName := *arguments.FileToCompile
		outputFilename := filenames.NsToQb(inputFileName)

		fmt.Printf("Compiling '%s'...\n", inputFileName)
		data, err := ioutil.ReadFile(inputFileName)
		check(err)

		sourceCode := NeverScript.NewSourceCode(string(data))
		byteCode, err := compilerService.Compile(sourceCode)
		check(err)

		err = ioutil.WriteFile(outputFilename, byteCode.GetBytes(), 0777)
		check(err)
		fmt.Printf("  Created '%s'.\n", outputFilename)

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
		flag.Usage()
	}
}

type commandLineArguments struct {
	FileToCompile   *string
	FileToDecompile *string
}

func parseCommandLineArguments() commandLineArguments {
	args := commandLineArguments{
		FileToCompile:   flag.String("c", "", "Specify a file to compile (.ns)."),
		FileToDecompile: flag.String("d", "", "Specify a file to decompile (.qb)."),
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
