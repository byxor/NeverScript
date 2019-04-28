package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/byxor/NeverScript/compiler"
)

const (
    banner = `
     __                    __           _       _   
  /\ \ \_____   _____ _ __/ _\ ___ _ __(_)_ __ | |_ 
 /  \/ / _ \ \ / / _ \ '__\ \ / __| '__| | '_ \| __|
/ /\  /  __/\ V /  __/ |  _\ \ (__| |  | | |_) | |_ 
\_\ \/ \___| \_/ \___|_|  \__/\___|_|  |_| .__/ \__|
                                         |_|        
`
)

func main() {
	arguments := parseCommandLineArguments()

	fmt.Println(banner)

	argumentsWereSupplied := false

	if *arguments.FileToCompile != "" {
		argumentsWereSupplied = true

		fmt.Printf("Compiling '%s'... ", *arguments.FileToCompile)

		data, err := ioutil.ReadFile(*arguments.FileToCompile)
		check(err)

		code := string(data)

		bytecode, err := compiler.Compile(code)
		check(err)

		fmt.Println("done")

		fmt.Println(bytecode)
	}

	if *arguments.FileToDecompile != "" {
		argumentsWereSupplied = true

		fmt.Printf("Compiling '%s'... ", *arguments.FileToDecompile)
		fmt.Print("(not implemented yet)")
		fmt.Println("done")
	}

	if !argumentsWereSupplied {
		flag.Usage()
	}
}

type commandLineArguments struct {
    FileToCompile *string
    FileToDecompile *string
}

func parseCommandLineArguments() commandLineArguments {
	args := commandLineArguments {
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
