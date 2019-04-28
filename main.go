package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/byxor/NeverScript/compiler"
)

func main() {
	arguments := parseCommandLineArguments()

	fmt.Println("Welcome to NeverScript")

	if *arguments.FileToCompile != "" {
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
		fmt.Printf("Compiling '%s'... ", *arguments.FileToDecompile)

		// check(err)

		fmt.Println("done")
	}

	fmt.Println("Goodbye!")
}

type commandLineArguments struct {
    FileToCompile *string
    FileToDecompile *string
}

func parseCommandLineArguments() commandLineArguments {
	args := commandLineArguments {
		FileToCompile:   flag.String("c", "", "A .ns file to compile."),
		FileToDecompile: flag.String("d", "", "A .qb file to decompile."),
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
