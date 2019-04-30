package compiler

import (
	"bytes"
	"fmt"
	"reflect"
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
				{"is_awesome = true;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0x01, 0x00, 0x00, 0x00,
					0x01, 0x00}},

				{"you_lost = false;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0x00, 0x00, 0x00, 0x00,
					0x01, 0x00}},
			})

			testThat("Integer variables can be declared", []testEntry{
				{"high_score = 0;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0x00, 0x00, 0x00, 0x00,
					0x01, 0x00}},

				{"high_score = 1;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0x01, 0x00, 0x00, 0x00,
					0x01, 0x00}},

				{"num_players = 255;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0xFF, 0x00, 0x00, 0x00,
					0x01, 0x00}},

				{"max_value = 4294967295;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0xFF, 0xFF, 0xFF, 0xFF,
					0x01, 0x00}},
			})

			testThat("Integer variables can be declared in a hexadecimal format", []testEntry{
				{"x_scale = 0x0;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0x00, 0x00, 0x00, 0x00,
					0x01, 0x00}},

				{"y_scale = 0xFFFFFFFF;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0xFF, 0xFF, 0xFF, 0xFF,
					0x01, 0x00}},

				{"z_scale = 0xBADBABE;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0xBE, 0xBA, 0xAD, 0x0B,
					0x01, 0x00}},
			})

			testThat("Integer variables can be declared in a binary format", []testEntry{
				{"larry = 0b0;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0x00, 0x00, 0x00, 0x00,
					0x01, 0x00}},

				{"silverstein = 0b1010;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0x0A, 0x00, 0x00, 0x00,
					0x01, 0x00}},

				{"max_value = 0b11111111111111111111111111111111;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0xFF, 0xFF, 0xFF, 0xFF,
					0x01, 0x00}},
			})

			testThat("Integer variables can be declared in an octal format", []testEntry{
				{"zero = 0o0;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0x00, 0x00, 0x00, 0x00,
					0x01, 0x00}},

				{"eight = 0o10;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0x08, 0x00, 0x00, 0x00,
					0x01, 0x00}},

				{"max_value = 0o37777777777;", []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x17, 0xFF, 0xFF, 0xFF, 0xFF,
					0x01, 0x00}},
			})

			testThat("String variables can be declared", []testEntry{
				{`empty = "";`, []byte{
					0x01, 0x16, any, any, any, any,
					0x07, 0x1B, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x01, 0x00}},

				{`name = "byxor";`, makeBytes(
					0x01, 0x16, any, any, any, any,
					0x07, 0x1B, 0x05, 0x00, 0x00, 0x00,
					"byxor", 0x00,
					0x01, 0x00)},

				{`weapon = "EML";`, makeBytes(
					0x01, 0x16, any, any, any, any,
					0x07, 0x1B, 0x03, 0x00, 0x00, 0x00,
					"EML", 0x00,
					0x01, 0x00)},
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

	test := func() {
		for _, entry := range entries {
			bytecode, err := compiler.Compile(entry.code)
			So(err, ShouldBeNil)

			replaceIrrelevantBytesFromFirstArgument(bytecode, entry.expectedBytecode)
			So(bytecode, shouldContainSubsequence, entry.expectedBytecode)
		}
	}

	Convey(description, test)
}

func replaceIrrelevantBytesFromFirstArgument(first, second []byte) {
	length := min(len(first), len(second))

	for i := 0; i < length; i++ {
		currentByteIsIrrelevant := second[i] == any

		if currentByteIsIrrelevant {
			first[i] = any
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

func makeBytes(elements ...interface{}) []byte {
	theBytes := make([]byte, 1024)
	size := 0

	for _, element := range elements {
		if theInt, ok := element.(int); ok {
			theByte := byte(theInt)
			theBytes[size] = theByte
			size++
			continue
		}

		if theString, ok := element.(string); ok {
			for _, theRune := range theString {
				theBytes[size] = byte(theRune)
				size++
			}
			continue
		}

		fmt.Println(reflect.TypeOf(element))
	}

	fmt.Println(theBytes[:size])

	return theBytes[:size]
}
