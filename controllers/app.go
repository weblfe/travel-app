package controllers

import "github.com/weblfe/travel-app/repositories"

type AppController struct {
		BaseController
}

func AppControllerOf() *AppController {
		return new(AppController)
}

// @router /app/config   [get]
func (this *AppController) GetGlobalConfig() {
		var driver = this.GetDriver()
		this.Send(repositories.NewAppRepository(this).GetConfig(driver))
}
