package test

import (
	"strings"
	"testing"

	"github.com/byxor/qbd/code"
	. "github.com/byxor/qbd/tokens"
	"github.com/stretchr/testify/assert"
)

const any = 0x00

func TestSyntax(t *testing.T) {
	entries := []struct {
		tokens   []Token
		expected string
	}{
		{[]Token{}, ""},
		{[]Token{{EndOfFile, nil}}, ""},

		// Line endings ------------------------------------------------------
		{[]Token{{EndOfLine, nil}}, ";"},
		{[]Token{{EndOfLine, nil}, {EndOfLine, nil}}, ";"},
		{[]Token{{EndOfLine, nil}, {EndOfLine, nil}, {EndOfLine, nil}}, ";"},

		// Integers ------------------------------------------------------
		{[]Token{{Integer, []byte{any, 0x00, 0x00, 0x00, 0x00}}}, "0"},
		{[]Token{{Integer, []byte{any, 0x01, 0x00, 0x00, 0x00}}}, "1"},
		{[]Token{{Integer, []byte{any, 0xFF, 0xFF, 0xFF, 0xFF}}}, "-1"},
		{[]Token{{Integer, []byte{any, 0x2A, 0x43, 0x0F, 0x00}}}, "1000234"},

		// Floats ------------------------------------------------------------
		{[]Token{{Float, []byte{any, 0x00, 0x00, 0xA0, 0x40}}}, "5.00f"},
		{[]Token{{Float, []byte{any, 0x33, 0x33, 0x03, 0x41}}}, "8.20f"},
		{[]Token{{Float, []byte{any, 0x33, 0xB3, 0x05, 0x43}}}, "133.70f"},

		// Pairs ------------------------------------------------------------
		{[]Token{{Pair, []byte{any, 0x66, 0x66, 0x66, 0x3F, 0x00, 0x00, 0x80, 0x3F}}}, "vec2<0.90f, 1.00f>"},
		{[]Token{{Pair, []byte{any, 0xCD, 0xCC, 0x4C, 0x3F, 0x00, 0x00, 0x80, 0x3F}}}, "vec2<0.80f, 1.00f>"},
		{[]Token{{Pair, []byte{any, 0x00, 0x00, 0xA0, 0x40, 0x00, 0x00, 0xA0, 0x40}}}, "vec2<5.00f, 5.00f>"},

		// Vectors ------------------------------------------------------------
		{[]Token{{Vector, []byte{any, 0x66, 0x66, 0x66, 0x3F, 0x00, 0x00, 0x01, 0x16, 0xCD, 0xCC, 0x4C, 0x3F}}}, "vec3<0.90f, 0.00f, 0.80f>"},
		{[]Token{{Vector, []byte{any, 0x00, 0x00, 0xC8, 0x42, 0x00, 0x00, 0x48, 0x42, 0xCD, 0xCC, 0x4C, 0x3F}}}, "vec3<100.00f, 50.00f, 0.80f>"},
		{[]Token{{Vector, []byte{any, 0x00, 0x00, 0x70, 0x41, 0x00, 0x00, 0x01, 0x16, 0x00, 0x00, 0x70, 0x41}}}, "vec3<15.00f, 0.00f, 15.00f>"},

		// Unknown Names ------------------------------------------------------
		{[]Token{{Name, []byte{any, 0x11, 0x22, 0x33, 0x44}}}, "%11223344%"},
		{[]Token{{Name, []byte{any, 0x12, 0x00, 0x00, 0x99}}}, "%12000099%"},
		{[]Token{{Name, []byte{any, 0xFF, 0xFF, 0xFF, 0xFF}}}, "%ffffffff%"},

		// Known Names ------------------------------------------------------
		{[]Token{
			{Name, []byte{any, 0x00, 0x00, 0x00, 0x00}},
			{NameTableEntry, []byte{
				any, 0x00, 0x00, 0x00, 0x00,
				0x47, 0x45, 0x54, 0x44, 0x4F, 0x57, 0x4E, 0x00},
			}},
			"GETDOWN",
		},

		{[]Token{
			{Name, []byte{any, 0xFF, 0xEE, 0x00, 0xCC}},
			{NameTableEntry, []byte{
				any, 0xFF, 0xEE, 0x00, 0xCC,
				0x54, 0x75, 0x72, 0x62, 0x6F, 0x54, 0x69, 0x6D, 0x65, 0x00},
			}},
			"TurboTime",
		},

		// Strings -------------------------------------------------------
		{[]Token{{String, []byte{
			any, 0x0D, 0x00, 0x00, 0x00, 0x4A, 0x6F,
			0x79, 0x20, 0x44, 0x69, 0x76, 0x69, 0x73,
			0x69, 0x6F, 0x6E, 0x00}}},
			"\"Joy Division\""},

		{[]Token{
			{String, []byte{any, 0x03, 0x00, 0x00, 0x00, 0x48, 0x69, 0x00}}},
			"\"Hi\""},

		// Addition  ------------------------------------------------------
		{[]Token{
			{Integer, []byte{any, 0x09, 0x00, 0x00, 0x00}},
			{Addition, nil},
			{Integer, []byte{any, 0x0A, 0x00, 0x00, 0x00}}},
			"9 + 10",
		},
		{[]Token{
			{Name, []byte{any, 0x09, 0x00, 0x00, 0x00}},
			{Addition, nil},
			{Name, []byte{any, 0x0A, 0x00, 0x00, 0x00}},
			{NameTableEntry, []byte{any, 0x0A, 0x00, 0x00, 0x00, 0x66, 0x6F, 0x6F, 0x00}}},
			"%09000000% + foo",
		},

		// Subtraction ------------------------------------------------------
		{[]Token{
			{EndOfLine, nil},
			{Integer, []byte{any, 0x00, 0xFF, 0xFF, 0xFF}},
			{Subtraction, nil},
			{Integer, []byte{any, 0x00, 0x00, 0x00, 0x00}}},
			"; -256 - 0",
		},
		{[]Token{
			{Name, []byte{any, 0x09, 0x00, 0x00, 0x00}},
			{Subtraction, nil},
			{LocalReference, nil},
			{Name, []byte{any, 0x0A, 0x00, 0x00, 0x00}}},
			"%09000000% - $%0a000000%",
		},

		// Division ------------------------------------------------------
		{[]Token{
			{Integer, []byte{any, 0xA0, 00, 00, 00}},
			{Division, nil},
			{Integer, []byte{any, 0x02, 00, 00, 00}}},
			"160 / 2",
		},
		{[]Token{
			{LocalReference, nil},
			{Name, []byte{any, 0xFF, 0xFF, 0xFF, 0xFF}},
			{Division, nil},
			{Integer, []byte{any, 0x20, 00, 00, 00}},
			{NameTableEntry, []byte{any, 0xFF, 0xFF, 0xFF, 0xFF, 0x61, 0x6E, 0x67, 0x6C, 0x65, 0x00}}},
			"$angle / 32",
		},

		// Multiplication ------------------------------------------------------
		{[]Token{
			{Integer, []byte{any, 0x02, 0x00, 0x00, 0x00}},
			{Multiplication, nil},
			{Integer, []byte{any, 0x03, 0x00, 0x00, 0x00}}},
			"2 * 3",
		},
		{[]Token{
			{LocalReference, nil},
			{Name, []byte{any, 0xFF, 0xFF, 0xFF, 0xFF}},
			{Multiplication, nil},
			{Integer, []byte{any, 0x20, 00, 00, 00}},
			{NameTableEntry, []byte{any, 0xFF, 0xFF, 0xFF, 0xFF, 0x61, 0x6E, 0x67, 0x6C, 0x65, 0x00}}},
			"$angle * 32",
		},

		// Assignment ------------------------------------------------------
		{[]Token{
			{Name, []byte{any, 0x00, 0x00, 0x00, 0x00}},
			{Assignment, nil},
			{Name, []byte{any, 0x11, 0x11, 0x11, 0x11}}},
			"%00000000% = %11111111%",
		},

		// Expressions ------------------------------------------------------
		{[]Token{
			{StartOfExpression, nil},
			{Integer, []byte{any, 0x09, 0x00, 0x00, 0x00}},
			{Addition, nil},
			{Integer, []byte{any, 0x0A, 0x00, 0x00, 0x00}},
			{EndOfExpression, nil}},
			"(9 + 10)",
		},

		// Arrays ------------------------------------------------------
		{[]Token{{StartOfArray, nil}, {EndOfArray, nil}}, "[]"},
		{[]Token{
			{StartOfArray, nil},
			{Integer, []byte{any, 0x12, 0x34, 0x56, 0x78}},
			{Integer, []byte{any, 0x00, 0x00, 0x00, 0x00}},
			{Integer, []byte{any, 0x0A, 0x00, 0x00, 0x00}},
			{Integer, []byte{any, 0xFF, 0xFF, 0xFF, 0xFF}},
			{StartOfExpression, nil},
			{Integer, []byte{any, 0x09, 0x00, 0x00, 0x00}},
			{Addition, nil},
			{Integer, []byte{any, 0x0A, 0x00, 0x00, 0x00}},
			{EndOfExpression, nil},
			{EndOfArray, nil}},
			"[2018915346 0 10 -1 (9 + 10)]",
		},
		{[]Token{
			{StartOfArray, nil},
			{StartOfArray, nil},
			{Name, []byte{any, 0xFF, 0x00, 0x00, 0xDD}},
			{EndOfArray, nil},
			{StartOfArray, nil},
			{Name, []byte{any, 0xBB, 0xEE, 0xEE, 0xFF}},
			{EndOfArray, nil},
			{EndOfArray, nil}},
			"[[%ff0000dd%] [%bbeeeeff%]]",
		},
		{[]Token{
			{EndOfLine, nil},
			{Name, []byte{any, 0xEF, 0xEF, 0xEF, 0xEF}},
			{Assignment, nil},
			{StartOfArray, nil},
			{StartOfArray, nil},
			{Name, []byte{any, 0xFF, 0x00, 0x00, 0xDD}},
			{EndOfArray, nil},
			{StartOfArray, nil},
			{Name, []byte{any, 0xBB, 0xEE, 0xEE, 0xFF}},
			{EndOfArray, nil},
			{EndOfArray, nil},
			{EndOfLine, nil},
			{Name, []byte{any, 0xEE, 0xEE, 0xEE, 0xEE}},
			{Assignment, nil},
			{StartOfArray, nil},
			{StartOfArray, nil},
			{Name, []byte{any, 0x11, 0x22, 0x33, 0x44}},
			{EndOfArray, nil},
			{StartOfArray, nil},
			{LocalReference, nil},
			{Name, []byte{any, 0x55, 0x66, 0x77, 0x88}},
			{EndOfArray, nil},
			{EndOfArray, nil},
			{EndOfLine, nil},
			{Name, []byte{any, 0x90, 0x80, 0x70, 0x60}},
			{Assignment, nil},
			{StartOfArray, nil},
			{StartOfArray, nil},
			{StartOfArray, nil},
			{StartOfArray, nil},
			{EndOfArray, nil},
			{EndOfArray, nil},
			{EndOfArray, nil},
			{EndOfArray, nil},
			{NameTableEntry, []byte{any, 0xEF, 0xEF, 0xEF, 0xEF, 0x66, 0x6F, 0x6F, 0x00}},
			{NameTableEntry, []byte{any, 0xEE, 0xEE, 0xEE, 0xEE, 0x62, 0x61, 0x72, 0x00}},
			{NameTableEntry, []byte{any, 0x90, 0x80, 0x70, 0x60, 0x62, 0x61, 0x7A, 0x00}}},
			lines(
				"; foo = [[%ff0000dd%] [%bbeeeeff%]]",
				"; bar = [[%11223344%] [$%55667788%]]",
				"; baz = [[[[]]]]",
			),
		},
	}

	for _, entry := range entries {
		code := code.GenerateUsing(entry.tokens)
		assert.Equal(t, entry.expected, code)
	}
}

func lines(lines ...string) string {
	return strings.Join(lines, "\n")
}
