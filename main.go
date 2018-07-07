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

	fmt.Println("DECOMPILING!")

	tokenChannel := make(chan tokens.Token)
	go tokens.Extract(tokenChannel, bytes)

	for {
		token, more := <-tokenChannel
		fmt.Println(more)
		if more {
			fmt.Println(token)
		} else {
			break
		}
	}

	fmt.Println("DONE!")
}
