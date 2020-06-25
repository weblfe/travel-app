package routers

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
}
