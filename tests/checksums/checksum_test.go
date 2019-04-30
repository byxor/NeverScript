package checksums

import (
	"github.com/byxor/NeverScript"
	"github.com/byxor/NeverScript/checksums"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var checksumService = checksums.NewService()

func TestChecksums(t *testing.T) {
	Convey("Checksums are correctly generated from identifiers", t, func() {
		data := []struct {
			identifier               string
			expectedChecksumContents uint32
		}{
			{"", 0x00},
			{"OffMeterTop", 0xD77B6FF9},
			{"offmetertop", 0xD77B6FF9},
			{"SetTrickScore", 0xCB3A8FD2},
			{"OneFootDarkSlide_range", 0x0750B702},
		}
		for _, entry := range data {
			actualChecksum := checksumService.GenerateFrom(entry.identifier)
			expectedChecksum := NeverScript.NewChecksum(entry.expectedChecksumContents)

			So(actualChecksum.IsEqualTo(expectedChecksum), ShouldBeTrue)
		}
	})
}

func TestLittleEndian(t *testing.T) {
	Convey("Checksums are correctly converted to little endian", t, func() {
		data := []struct {
			checksumContents uint32
			expectedBytes    []byte
		}{
			{0xD77B6FF9, []byte{0xF9, 0x6F, 0x7B, 0xD7}},
			{0x12345678, []byte{0x78, 0x56, 0x34, 0x12}},
		}
		for _, entry := range data {
			checksum := NeverScript.NewChecksum(entry.checksumContents)
			actualBytes := checksumService.EncodeAsLittleEndian(checksum)
			So(actualBytes, ShouldResemble, entry.expectedBytes)
		}
	})
}
