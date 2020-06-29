package test

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/config/env"
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/libs"
		"testing"
)

func TestVariableParse(t *testing.T) {
		var (
				os      = "mac"
				appName = beego.AppConfig.String("appname")
				str     = "${app_name}/${os}/${def|sdf}/123213"
				str1     = appName+`/mac/sdf/123213`
		)
		env.Set("app_name", appName)
		env.Set("os", os)
		Convey("Test Var parse",t, func() {
				So(libs.VariableParse(str),ShouldEqual,str1)
		})

}
