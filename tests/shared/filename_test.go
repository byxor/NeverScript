package shared

import (
	"github.com/byxor/NeverScript/shared/filenames"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConvertingQbToNs(t *testing.T) {
	Convey(".qb filenames are converted .qb filenames", t, func() {
    	entries := []testEntry{
    		{"a.qb", "a.ns"},
    		{"foo.qb", "foo.ns"},
			{"~/foo/bar/baz.qb", "~/foo/bar/baz.ns"},
    	}
    	for _, entry := range entries {
			output := filenames.QbToNs(entry.input)
    		So(output, ShouldEqual, entry.expectedOutput)
    	}
	})
}

type testEntry struct {
	input          string
	expectedOutput string
}
