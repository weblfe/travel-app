package test

import (
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
		"testing"
)

func TestSmsCodeServiceSendCode(t *testing.T)  {
		var (
				typ = "register"
			mobile = "13112260988"
			extras = map[string]string{}
		)
		common.SmsDebugOn()
		Convey("Test Sms Code",t, func() {
				code,err:=services.SmsCodeServiceOf().SendCode(mobile,typ, extras)
				So(err,ShouldBeNil)
				So(code,ShouldNotBeEmpty)
				res:=services.SmsCodeServiceOf().Verify(mobile,code,typ)
				So(res,ShouldBeTrue)
		})

}
