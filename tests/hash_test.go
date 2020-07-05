package test

import (
		"fmt"
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"testing"
)

func TestHashCode(t *testing.T) {
		var hash string
		Convey("hash Test", t, func() {
				hash = libs.HashCode("231234")
				So(hash, ShouldNotBeEmpty)
				var obj = models.UserModelOf()
				hash1 := libs.HashCode(*obj)
				hash2 := libs.HashCode(*models.UserModelOf())
				hash3 := libs.HashCode(*obj)
			//	fmt.Println(hash,hash2,hash1,hash3)
				So(hash1, ShouldNotEqual, hash2)
				So(hash1, ShouldNotEqual, hash3)
		})
}

func BenchmarkHashCode(b *testing.B) {
		for i := 0; i < b.N; i++ {
				obj := *models.UserModelOf()
				_ = libs.HashCode(obj)
				_ = libs.HashCode(*models.UserModelOf())
		}
}

func ExampleHashCode() {
		var hash = libs.HashCode("231234")
		fmt.Println(hash)
		obj := *models.UserModelOf()
		hash1 := libs.HashCode(obj)
		fmt.Println(hash1)
		hash2 := libs.HashCode(*models.UserModelOf())
		fmt.Println(hash2)
}
