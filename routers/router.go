package routers

// @APIVersion 1.0.0
// @Title beego Test API
// @Description api document
// @Contact weblinuxgame@g126.com
// @TermsOfServiceUrl http://api.word-server.com/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/controllers"
)

func init() {
		// 默认 /
		beego.Include(controllers.IndexControllerOf())
		// 用户模块
		beego.Include(controllers.UserControllerOf())
		// 作品模块
		beego.Include(controllers.PostsControllerOf())
		// 验证码模块
		beego.Include(controllers.CaptchaControllerOf())
		// 消息模块
		beego.Include(controllers.MessageControllerOf())
		// 附件模块
		beego.Include(controllers.AttachmentControllerOf())
		// 点赞模块
		beego.Include(controllers.ThumbsUpControllerOf())
}
