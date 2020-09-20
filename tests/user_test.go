package test

import (
		"fmt"
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/libs"
		"testing"
)

func TestRandomNickName(t *testing.T) {
		Convey("Test random nickname", t, func() {
				var name = libs.RandomNickName(20)
				fmt.Println(name)
				So(len(name) == 16, ShouldBeTrue)
		})
}

func BenchmarkRandomNickName(b *testing.B) {
		for i := 0; i < b.N; i++ {
				name := libs.RandomNickName(20)
				fmt.Println(name)
		}
}
