package main

import (
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
	go tokens.Extract(tokenChannel, bytes)

	for token := range tokenChannel {
		color.Green(token.String())
	}

	fmt.Println("Stopped decompilation.")
}
