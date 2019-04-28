package shared

import (
	"github.com/byxor/NeverScript/compiler/checksums"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestChecksums(t *testing.T) {
	Convey("Checksums are generated correctly", t, func() {
		data := []struct {
			name     string
			checksum uint32
		}{
			{"", 0x00},
			{"OffMeterTop", 0xD77B6FF9},
			{"SetTrickScore", 0xCB3A8FD2},
			{"OneFootDarkSlide_range", 0x0750B702},
		}
		for _, entry := range data {
			result := checksums.Generate(entry.name)
			So(result, ShouldResemble, entry.checksum)
		}
	})
}

func TestLittleEndian(t *testing.T) {
	Convey("Checksums are correctly converted to little endian", t, func() {
		data := []struct {
			checksum     uint32
			littleEndian []byte
		}{
			{0xD77B6FF9, []byte{0xF9, 0x6F, 0x7B, 0xD7}},
			{0x12345678, []byte{0x78, 0x56, 0x34, 0x12}},
		}
		for _, entry := range data {
			result := checksums.LittleEndian(entry.checksum)
			So(result, ShouldResemble, entry.littleEndian)
		}
	})
}
