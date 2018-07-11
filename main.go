package main

import (
	"fmt"
	"github.com/byxor/qbd/tokens"
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

	for range tokenChannel {
	}

	fmt.Println("Stopped decompilation.")
}
