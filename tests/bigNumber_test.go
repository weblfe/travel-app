package test

import (
		"fmt"
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/repositories"
		"testing"
)

func TestBigNumber(t *testing.T) {
		var numbers = []map[string]int64{
				{"1": 1, "0": 0, "11": 11, "123": 123, "0.99k": 999},
				{"1k": 1000, "1.2k": 1200, "9k": 9 * repositories.BigNumberK, "1w": repositories.BigNumberW+1,"10w":int64(repositories.BigNumberW*10)},
				{"1.23k": 1230, "1kw": repositories.BigNumberKW, "1ww": repositories.BigNumberWW},
		}
		Convey("Test big Number ", t, func() {
				for _, nums := range numbers {
						for unit, num := range nums {
								text := repositories.DecorateNumberToText(num)
								fmt.Println(text)
								So(text == unit, ShouldBeTrue)
						}
				}
		})
}
