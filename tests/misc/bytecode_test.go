package misc

import (
	"github.com/byxor/NeverScript"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestByteCode(t *testing.T) {
	byteCode := NeverScript.NewByteCode([]byte{0x00, 0x01, 0x02, 0x03})

	Convey("When getting a slice of ByteCode", t, func() {
		Convey("Errors are raised when", func() {
			Convey("The start index is negative", func() {
				_, err := byteCode.GetSlice(-1, 1)
				So(err, ShouldEqual, NeverScript.SliceIndexOutOfRange)
			})

			Convey("The end index is negative", func() {
				_, err := byteCode.GetSlice(0, -1)
				So(err, ShouldEqual, NeverScript.SliceIndexOutOfRange)
			})

			Convey("The end index is below the start index", func() {
				_, err := byteCode.GetSlice(2, 1)
				So(err, ShouldEqual, NeverScript.SliceIndexOutOfRange)
			})

			Convey("The end index is above the ByteCode length", func() {
				_, err := byteCode.GetSlice(0, 5)
				So(err, ShouldEqual, NeverScript.SliceIndexOutOfRange)
			})
		})
	})
}
