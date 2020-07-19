package test

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/transforms"
		"reflect"
		"testing"
		"time"
)

func TestIsEmail(t *testing.T) {
		var emails = getEmails()
		Convey("Test Regexp Email", t, func() {
				for _, it := range emails {
						So(libs.IsEmail(it.Key), ShouldEqual, it.Value)
				}
		})
}

func TestIsReflect(t *testing.T) {
		var emails []string
		emails = append(emails, "34@qq.com", "hello@163.com")
		getValue := reflect.ValueOf(emails)
		getType := reflect.TypeOf(emails)
		fmt.Println(getType.Kind() == reflect.Array, getValue.Kind() == reflect.Slice)
		// fmt.Println(getType.Len())
		fmt.Println(getValue.Len())
}

func TestFilter(t *testing.T) {
		var (
				filters []func(m beego.M) beego.M
				mapper  = models.NewAttachment()
		)
		filters = append(filters, func(m beego.M) beego.M {
				for key, v := range m {
						obj, ok := v.(map[string]interface{})
						if ok && len(obj) == 0 {
								delete(m, key)
								continue
						}
						getValue := reflect.ValueOf(v)
						if getValue.Kind() == reflect.Map && getValue.Len() == 0 {
								delete(m, key)
								continue
						}
						obj, ok = v.(beego.M)
						if ok && len(obj) == 0 {
								delete(m, key)
								continue
						}
						obj, ok = v.(bson.M)
						if ok && len(obj) == 0 {
								delete(m, key)
						}
				}
				return m
		})

		filters = append(filters, func(m beego.M) beego.M {
				for k, v := range m {
						if v == "" || v == nil {
								delete(m, k)
								continue
						}
						getValue := reflect.ValueOf(v)
						kindName := getValue.Kind()
						if kindName == reflect.Array && getValue.Len() <= 0 {
								delete(m, k)
								continue
						}
						if getValue.IsZero() {
								delete(m, k)
								continue
						}
						if t, ok := v.(time.Time); ok {
								if t.IsZero() {
										delete(m, k)
								}
						}

				}
				return m
		})

		m:= mapper.M(transforms.FilterWrapper(filters...))
		fmt.Printf("%v\n", mapper)
		fmt.Printf("%v\n", m)
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
