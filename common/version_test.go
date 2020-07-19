package common

import (
		. "github.com/smartystreets/goconvey/convey"
		"testing"
)

func TestVersion_Check(t *testing.T) {
		var (
				v1 = "1.9.0"
				v2 = "1.0.10"
				v3 = "1.0"
				v4 = "0.0.1"
		)

		Convey("test version check", t, func() {
				So(Version(v1).Check(v1, "=="), ShouldBeTrue)
				So(Version(v2).Check(v1, "<"), ShouldBeTrue)
				So(Version(v3).Check(v1, ">"), ShouldBeFalse)
				So(Version(v4).Check(v1, "!="), ShouldBeTrue)
				So(Version(v4).Check(v1, "<"), ShouldBeTrue)
		})
}

func BenchmarkVersion_Check(b *testing.B) {
		var (
				v1 = "1.9.0"
			/*	v2 = "1.0.10"
				v3 = "1.0"
				v4 = "0.0.1"*/
		)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
				Version(v1).Check(v1, "==")
		}
		b.StopTimer()
}
