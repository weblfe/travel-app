package test

import (
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/libs"
		"testing"
)

func TestIsEmail(t *testing.T) {
		var emails = getEmails()
		Convey("Test Regexp Email", t, func() {
				for _, it := range emails {
						So(libs.IsEmail(it.Key), ShouldEqual, it.Value)
				}
		})
}

func TestIsMobile(t *testing.T) {
		var mobiles = getMobiles()
		Convey("Test Regexp Mobile", t, func() {
				for _, it := range mobiles {
						m := libs.IsCnMobile(it.Key)
						So(m, ShouldEqual, it.Value)
				}
		})
}

func getMobiles() []struct {
		Key   string
		Value bool
} {
		return []struct {
				Key   string
				Value bool
		}{
				{Key: "13112260977", Value: true},
				{Key: "13112260970", Value: true},
				{Key: "1311226097", Value: false},
				{Key: "23112260977", Value: false},
				{Key: "131122609xx", Value: false},
				{Key: "12112260977", Value: true},
				{Key: "10000000000", Value: false},
		}
}

func getEmails() []struct {
		Key   string
		Value bool
} {
		return []struct {
				Key   string
				Value bool
		}{
				{Key: "99468@qq.com", Value: true},
				{Key: "weblinuxgame@126.com", Value: true},
				{Key: "weblinuxgame@baidu...se", Value: false},
				{Key: "9@11121233", Value: false},
				{Key: "#sdfsf@111.com", Value: false},
				{Key: "weblinuxgame.com", Value: false},
				{Key: "999@@@@.com", Value: false},
		}
}

func BenchmarkIsMobile(b *testing.B) {
		for i := 0; i < b.N; i++ {
				for _, it := range getMobiles() {
						libs.IsCnMobile(it.Key)
				}
		}
}

func BenchmarkRandomWords(b *testing.B) {
		for i := 0; i < b.N; i++ {
				for i := 0; i < 10; i++ {
						_ = libs.RandomAnyWord(i + 1)
						_ = libs.RandomWord(i + 1)
				}
		}
}

func TestRandomWords(t *testing.T) {
		Convey("Test Random Word", t, func() {
				for i := 0; i < 10; i++ {
						word := libs.RandomAnyWord(i + 1)
						So(len([]rune(word)), ShouldEqual, i+1)
				}
		})
}
