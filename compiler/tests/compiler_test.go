package tests

import (
	"github.com/byxor/NeverScript/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCompilation(t *testing.T) {
	Convey("When compiling QB code", t, func() {
		Convey("We get bytecode", func() {
			data := []struct {
				code     string
				bytecode []byte
			}{
				{"", []byte{0x00}},
				{";", []byte{0x01, 0x00}},
				{";;", []byte{0x01, 0x01, 0x00}},
				{";;;", []byte{0x01, 0x01, 0x01, 0x00}},

				{";  ;\t;\n;", []byte{0x01, 0x01, 0x01, 0x01, 0x00}},

				{";foo = 10", []byte{
					0x01,
					0x16, 0x49, 0x73, 0x18, 0x61,
					0x07, 0x0A, 0x00, 0x00, 0x00,
					0x01,
					0x2B, 0x49, 0x73, 0x18, 0x61,
					0x66, 0x6F, 0x6F, 0x00,
					0x00,
				}},
			}
			for _, entry := range data {
				result, err := qbc.Compile(entry.code)
				So(err, ShouldBeNil)
				So(result, ShouldResemble, entry.bytecode)
			}
		})
	})
}
