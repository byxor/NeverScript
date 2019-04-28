package compiler

import (
	"bytes"
	"fmt"
	"github.com/byxor/NeverScript/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const any = 0xF4

func TestCompilation(t *testing.T) {
	Convey("When compiling NeverScript code", t, func() {
		Convey("We get QB bytecode", func() {
			testThat("Files end with 0x00", []testEntry{
				{"", []byte{0x00}},
			})

			testThat("Semicolons insert 0x01", []testEntry{
				{";", []byte{0x01, 0x00}},
				{";;", []byte{0x01, 0x01, 0x00}},
				{";;;", []byte{0x01, 0x01, 0x01, 0x00}},
			})

			testThat("Whitespace is ignored", []testEntry{
				{";  ;\t;\n;", []byte{0x01, 0x01, 0x01, 0x01, 0x00}},
			})

			testThat("Boolean variables can be declared", []testEntry{
				{"is_awesome = true;", []byte{0x01, 0x16, any, any, any, any,
                                              0x07, 0x17, 0x01, 0x00, 0x00, 0x00,
					                          0x01, 0x00}},

				{"you_lost = false;", []byte{0x01, 0x16, any, any, any, any,
					                         0x07, 0x17, 0x00, 0x00, 0x00, 0x00,
					                         0x01, 0x00}},
			})
		})
	})
}

type testEntry struct {
	code             string
	expectedBytecode []byte
}

func testThat(someRequirementIsMet string, entries []testEntry) {
	description := fmt.Sprintf("Test: %s", someRequirementIsMet)

	Convey(description, func() {
		for _, entry := range entries {
			result, err := compiler.Compile(entry.code)

			So(err, ShouldBeNil)

			replaceIrrelevantBytesFromFirstArgument(result, entry.expectedBytecode)
			So(result, shouldContainSubsequence, entry.expectedBytecode)
		}
	})
}

func shouldContainSubsequence(actual interface{}, expected ...interface{}) string {
	sequence := actual.([]byte)
	subsequence := expected[0].([]byte)

	sequenceNotFound := fmt.Sprintf(
		"%s\nSequence:    %v\nSubsequence: %v\n",
		"Sequence doesn't contain expected subsequence.",
		sequence,
		subsequence,
	)

	if len(subsequence) > len(sequence) {
		return sequenceNotFound
	}

	for i := 0; i < len(sequence)-len(subsequence)+1; i++ {
		sliceOfSequence := sequence[i : i+len(subsequence)]

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
