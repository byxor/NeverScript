package compiler

import (
	"github.com/byxor/NeverScript/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"bytes"
	"fmt"
	"testing"
)

const any = 0xF4

type testData struct {
	code             string
	expectedBytecode []byte
}

func TestCompilation(t *testing.T) {
	Convey("When compiling NeverScript code", t, func() {
		Convey("We get QB bytecode", func() {
			data := []testData{
				// File ends with 0x00
				{"", []byte{0x00}},

				// Semicolons insert 0x01
				{";", []byte{0x01, 0x00}},
				{";;", []byte{0x01, 0x01, 0x00}},
				{";;;", []byte{0x01, 0x01, 0x01, 0x00}},

				// Whitespace is ignored
				{";  ;\t;\n;", []byte{0x01, 0x01, 0x01, 0x01, 0x00}},

				// Boolean variables can be declared
				{"is_awesome = true;", []byte{0x01, 0x16, any, any, any, any,
                                              0x07, 0x17, 0x01, 0x00, 0x00, 0x00,
				                              0x01, 0x00}},

				{"is_winner = false;", []byte{0x01, 0x16, any, any, any, any,
                                              0x07, 0x17, 0x00, 0x00, 0x00, 0x00,
				                              0x01, 0x00}},
			}
			for _, entry := range data {
				result, err := compiler.Compile(entry.code)

				So(err, ShouldBeNil)

				replaceIrrelevantBytesFromFirstArgument(result, entry.expectedBytecode)
				So(result, shouldContainSubsequence, entry.expectedBytecode)
			}
		})
	})
}

func shouldContainSubsequence(actual interface{}, expected ...interface{}) string {
	sequence := actual.([]byte)
	subsequence := expected[0].([]byte)

	fmt.Printf("Checking if %v\ncontains    %v\n", sequence, subsequence)

	sequenceNotFound := fmt.Sprintf(
		"%s\nSequence:    %v\nSubsequence: %v\n",
		"Sequence doesn't contain expected subsequence.",
		sequence,
		subsequence,
	)

	if len(subsequence) > len(sequence) {
		return sequenceNotFound
	}

	for i := 0; i < len(sequence) - len(subsequence) + 1; i++ {
		sliceOfSequence := sequence[i:i+len(subsequence)]
		fmt.Printf(".     Comparing %v with %v\n", sliceOfSequence, subsequence)
		if bytes.Equal(sliceOfSequence, subsequence) {
			sequenceFound := ""
			return sequenceFound
		}
	}

	return sequenceNotFound
}

func replaceIrrelevantBytesFromFirstArgument(actualBytes []byte, expectedBytes []byte) {
	length := min(len(actualBytes), len(expectedBytes))

	for i := 0; i < length; i++ {
		if expectedBytes[i] == any {
			actualBytes[i] = any
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
