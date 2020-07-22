package controllers

type AppController struct {
		BaseController
}

func AppControllerOf() *AppController  {
		return new(AppController)
}

// @router /app/config   [get]
func (this *AppController)GetGlobalConfig()  {
	this.Send(nil)
}