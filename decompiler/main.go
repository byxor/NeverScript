package main

import (
	"fmt"
	"github.com/byxor/NeverScript/decompiler/code"
	"github.com/byxor/NeverScript/decompiler/tokens"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	fileName := os.Args[1]

	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	tokenChannel := make(chan tokens.Token)
	go tokens.ExtractAll(tokenChannel, bytes)

	tokens := []tokens.Token{}
	for token := range tokenChannel {
		tokens = append(tokens, token)
	}

	code := code.GenerateUsing(tokens)

	fmt.Println(code)
}
