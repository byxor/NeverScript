package misc

import (
	"github.com/byxor/NeverScript"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConvertingQbToNs(t *testing.T) {
	Convey(".qb filenames are converted .ns filenames", t, func() {
		entries := []testEntry{
			{"a.qb", "a.ns"},
			{"foo.qb", "foo.ns"},
			{"~/foo/bar/baz.qb", "~/foo/bar/baz.ns"},
		}
		for _, entry := range entries {
			output := NeverScript.QbToNs(entry.input)
			So(output, ShouldEqual, entry.expectedOutput)
		}
	})
}

func TestConvertingNsToQb(t *testing.T) {
	Convey(".ns filenames are converted .qb filenames", t, func() {
		entries := []testEntry{
			{"a.ns", "a.qb"},
			{"foo.ns", "foo.qb"},
			{"~/foo/bar/baz.ns", "~/foo/bar/baz.qb"},
		}
		for _, entry := range entries {
			output := NeverScript.NsToQb(entry.input)
			So(output, ShouldEqual, entry.expectedOutput)
		}
	})
}

type testEntry struct {
	input          string
	expectedOutput string
}
