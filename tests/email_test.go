package test

import (
		. "github.com/smartystreets/goconvey/convey"
		"github.com/weblfe/travel-app/services"
		"os"
		"testing"
)

func TestEmail(t *testing.T) {
		var email = "15975798646@163.com"
		Convey("Test email sends", t, func() {
				request := services.NewEmailRequest()
				request.Set("to", email)
				request.Set("content", `<div>hello world!</div>`)
				request.Set("isHtml", "true")
				request.AddFile(os.Getenv("GOPATH") + "/src/github.com/weblfe/travel-app" + "/" + "Dockerfile")
				request.Set(services.EmailSubject, "测试邮件")
				request.Set(services.EmailFrom, "Jordan Wright <weblinuxgame@126.com>")
				do := services.EmailServiceOf().Sends(request)
				So(do, ShouldBeNil)
		})

}
