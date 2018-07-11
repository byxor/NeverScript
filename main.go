package main

import (
	"encoding/hex"
	"fmt"
	"github.com/byxor/qbd/tokens"
	"github.com/fatih/color"
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

	fmt.Println("Starting decompilation...")

	tokenChannel := make(chan tokens.Token)
	go tokens.ExtractAll(tokenChannel, bytes)

	for token := range tokenChannel {
		type colorFunction func(format string, a ...interface{})
		var displayType, displayChunk colorFunction

		if token.Type == tokens.Invalid {
			displayType = color.Red
			displayChunk = color.Red
		} else {
			displayType = color.Green
			displayChunk = color.White
		}

		displayType(token.Type.String())
		displayChunk(hex.Dump(token.Chunk) + "\n")
	}

	fmt.Println("Stopped decompilation.")
}
