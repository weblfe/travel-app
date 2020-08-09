package test

import (
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/plugins"
		"testing"
)

func TestQrcode(t *testing.T) {
		var (
				content = "id=app&code=test"
				qrcode  = plugins.GetQrcode()
			//	root    = "/Users/hewei/workpaces/coding/golang/src/github.com/weblfe/travel-app/static/storage/2020-8-9"
		)
		var (
				option = &plugins.Options{
						Size:            256,
						Level: 2,
						Auto: true,
					//	BackgroundColor: color.RGBA{R: 0x33, G: 0x33, B: 0x66, A: 0xff},
					//	ForegroundColor: color.RGBA{R: 0xef, G: 0xef, B: 0xef, A: 0xff},
				}
				err error
		)
		option.Init()
		err = qrcode.Save(content, option)
		Convey("test Qrcode err", t, func() {
				So(err == nil, ShouldBeTrue)
				m, err := qrcode.Decode(option.FileName)
				So(err == nil, ShouldBeTrue)
				So(m != nil, ShouldBeTrue)
			 //	fmt.Println(m.Content)
				So(m.Content == content, ShouldBeTrue)
		})
}
