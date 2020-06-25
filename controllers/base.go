package controllers

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
)

type BaseController struct {
		beego.Controller
}

func (this *BaseController) Send(json common.ResponseJson) {
		this.Data["json"] = json
		this.ServeJSON()
}
