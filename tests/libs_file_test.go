package test

import (
		"fmt"
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/libs"
		"math"
		"regexp"
		"testing"
)

func TestFileSize(t *testing.T) {
		var (
				size     int64
				fileSize libs.FileSize
		)
		size = int64(math.Pow(1024, 4))
		str := libs.FormatFileSize(size)
		tran := fileSize.Parse(str)
		Convey("Test file size", t, func() {
				So(str, ShouldNotBeEmpty)
				So(size, ShouldEqual, tran)
				size = size + 100
				str = libs.FormatFileSize(size)
				tran = fileSize.Parse(str)
				So(str, ShouldNotBeEmpty)
				So(tran-size <= 500, ShouldEqual, true)
		})
}

// E11000 duplicate key error collection: travel.users index: email_1 dup key: { email: \"994685561@qq.com\" }

func TestMatcher(t *testing.T) {
		var (
				str  = `E11000 duplicate key error collection: travel.users index: email_1 dup key: { email: "994685561@qq.com" }`
				reg  = regexp.MustCompile(`.+ dup key: (.+)`)
				regs = regexp.MustCompile(`.+ dup keys: (.+)`)
		)
		arr := reg.FindAllStringSubmatch(str,-1)
		fmt.Println(len(arr[0]))
		for _, it := range arr {
				for _,v:= range it {
						fmt.Println(v)
				}
		}
		arrAll := regs.FindAllString(str, -1)
		for _, it := range arrAll {
				fmt.Println(it)
		}
}
