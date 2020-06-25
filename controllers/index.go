package controllers

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
)

type IndexController struct {
		BaseController
}

func IndexControllerOf() *IndexController {
		return new(IndexController)
}

// @router /  [get]
// @router /  [post]
// @router /  [delete]
// @router /  [put]
func (this *IndexController) Index() {
		this.Send(common.NewResponse(beego.M{}, 0, common.Success))
}
