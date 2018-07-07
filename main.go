package main

import (
	"fmt"
	"github.com/byxor/qbd/token"
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

	fmt.Println("DECOMPILING!")

	tokens := make(chan token.Token)
	go token.GetTokens(tokens, bytes)

	for {
		token, more := <-tokens
		fmt.Println(more)
		if more {
			fmt.Println(token)
		} else {
			break
		}
	}

	fmt.Println("DONE!")
}
