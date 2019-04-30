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
