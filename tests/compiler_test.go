package tests

import (
	"fmt"
	"github.com/byxor/NeverScript"
	"github.com/byxor/NeverScript/compiler"
	"github.com/byxor/NeverScript/checksums"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"encoding/hex"
)

var (
	checksumService = checksums.NewService()
	compilerService = compiler.NewService(checksumService)
	any = byte(0xF4)
)

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
	sourceCodeContent       string
	expectedByteCodeContent []byte
}

func testThat(someRequirementIsMet string, entries []testEntry) {
	description := fmt.Sprintf("Test: %s", someRequirementIsMet)

	test := func() {
		for _, entry := range entries {
			sourceCode := NeverScript.NewSourceCode(entry.sourceCodeContent)
			expectedByteCode := NeverScript.NewByteCode(entry.expectedByteCodeContent)

			actualByteCode, err := compilerService.Compile(sourceCode)
			So(err, ShouldBeNil)

			So(actualByteCode, shouldContainSubsequence, expectedByteCode, any)
		}
	}

	Convey(description, test)
}

func shouldContainSubsequence(actual interface{}, expected ...interface{}) string {
	sequence, ok := actual.(NeverScript.ByteCode)
	if !ok {
		return "Couldn't cast 'sequence' to ByteCode"
	}

	subsequence, ok := expected[0].(NeverScript.ByteCode)
	if !ok {
		return "Couldn't cast 'subsequence' to ByteCode"
	}

	// temp, ok := expected[1].(byte)
	// if !ok {
	// 	return "Couldn't cast 'temp' to int"
	// }

	// byteToIgnore := byte(temp)

	sequenceNotFound := fmt.Sprintf(
		"%s\nSequence:    %v\nSubsequence: %v\n",
		"Sequence doesn't contain expected subsequence.",
		hex.Dump(sequence.ToBytes()),
		hex.Dump(subsequence.ToBytes()),
	)

	// containsSubsequence, err := sequence.Contains_IgnoreByte(subsequence, byteToIgnore)
	containsSubsequence, err := sequence.Contains(subsequence)
	if err != nil {
		return err.Error()
	}

	if !containsSubsequence {
		return sequenceNotFound
	}

	return ""
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
	}

	return theBytes[:size]
}

