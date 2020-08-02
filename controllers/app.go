package controllers

import "github.com/weblfe/travel-app/repositories"

type AppController struct {
		BaseController
}

func AppControllerOf() *AppController {
		return new(AppController)
}

// 获取应用相关配置
// @router /app/config   [get]
func (this *AppController) GetGlobalConfig() {
		var driver = this.GetDriver()
		this.Send(repositories.NewAppRepository(this).GetConfig(driver))
}

//  申请 ｜ 提交举报
// @router /app/apply  [post]
func (this *AppController) CommitApply() {
		this.Send(repositories.NewAppRepository(this).Apply())
}

