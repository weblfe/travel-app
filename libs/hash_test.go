package libs

import (
		"fmt"
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/models"
		"testing"
)

func TestHashCode(t *testing.T) {
		var hash string
		Convey("hash Test", t, func() {
				hash = HashCode("231234")
				So(hash, ShouldNotBeEmpty)
				var obj = models.UserModelOf()
				hash1 := HashCode(*obj)
				hash2 := HashCode(*models.UserModelOf())
				hash3 := HashCode(*obj)
			//	fmt.Println(hash,hash2,hash1,hash3)
				So(hash1, ShouldNotEqual, hash2)
				So(hash1, ShouldNotEqual, hash3)
		})
}

func BenchmarkHashCode(b *testing.B) {
		for i := 0; i < b.N; i++ {
				obj := *models.UserModelOf()
				_ = HashCode(obj)
				_ = HashCode(*models.UserModelOf())
		}
}

func ExampleHashCode() {
		var hash = HashCode("231234")
		fmt.Println(hash)
		obj := *models.UserModelOf()
		hash1 := HashCode(obj)
		fmt.Println(hash1)
		hash2 := HashCode(*models.UserModelOf())
		fmt.Println(hash2)
}
