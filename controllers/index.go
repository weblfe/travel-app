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

// @router /  [*]
func (this *IndexController) Index() {
		this.Send(common.NewResponse(beego.M{}, 0, common.Success))
}


// @router /app/about [get]
func (this *IndexController) GetAbout()  {
		this.View("about.tpl")
}

// @router /app/contactUs [get]
func (this *IndexController) GetContactUs()  {
		this.View("contactUs.tpl")
}

// @router /app/privacy [get]
func (this *IndexController) GetPrivacy()  {
		this.View("privacy.tpl")
}

// @router /app/agreement [get]
func (this *IndexController) GetAgreement()  {
		this.View("agreement.tpl")
}