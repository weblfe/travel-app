package test

import (
		"fmt"
		"github.com/weblfe/travel-app/models"
		_ "github.com/weblfe/travel-app/routers"
		"testing"

		. "github.com/smartystreets/goconvey/convey"
)

// TestBeego is a sample to run an endpoint test
func TestWord(t *testing.T) {
		var test = []struct {
				Text string
				Ok   bool
		}{
				{"AV做爱❤️", true},
				{"你妈逼的", true},
				{"操你老母", true},
				{"环球网报道】据香港“东网”等媒体2日报道，香港卫生署卫生防护中心今天（2日）表示，香港新增115例新冠肺炎确诊病例，均为本地确诊病例。这已是香港确诊病例连续第12日单日破百，截至目前，香港确诊病例累计3511例。", false},
				{"今天好开心", false},
		}
		models.SetProfile("db_name", "travel")
		Convey("", t, func() {
				for _, it := range test {
						newTxt := models.GetDfaInstance().ChangeSensitiveWords(it.Text)
						fmt.Println(it.Text, newTxt)
						if it.Ok {
								So(it.Text != newTxt, ShouldBeTrue)
						} else {
								So(it.Text == newTxt, ShouldBeTrue)
						}
				}
		})
}

func BenchmarkWord(b *testing.B) {
		var test = []struct {
				Text string
				Ok   bool
		}{
				/*{"AV做爱❤️", true},
				{"你妈逼的", true},
				{"操你老母", true},*/
				{"环球网报道】据香港“东网”等媒体2日报道，香港卫生署卫生防护中心今天（2日）表示，香港新增115例新冠肺炎确诊病例，均为本地确诊病例。这已是香港确诊病例连续第12日单日破百，截至目前，香港确诊病例累计3511例。傻逼，你做爱", false},
				/*{"今天好开心", false},*/
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
				for _, it := range test {
						_ = models.GetDfaInstance().ChangeSensitiveWords(it.Text)
						// fmt.Println(it.Text, newTxt)
				}
		}
}
