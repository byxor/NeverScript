package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}

func Decompile(bytes []byte) string {
	if len(bytes) == 0 {
		return ""
	} else {
		switch bytes[0] {
		case 0x01:
			return ";"
		case 0x16:
			return fmt.Sprintf("#%08x", binary.LittleEndian.Uint32(bytes[1:]))
		case 0x17:
			return fmt.Sprint(int32(binary.LittleEndian.Uint32(bytes[1:])))
		}
		return ""
	}
}
