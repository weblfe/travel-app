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

// Index
// @router /  [get]
// @router /  [post]
func (this *IndexController) Index() {
		this.Send(common.NewResponse(beego.M{}, 0, common.Success))
}

// GetAbout
// @router /app/about [get]
func (this *IndexController) GetAbout()  {
		this.View("about.tpl")
}

// GetContactUs
// @router /app/contactUs [get]
func (this *IndexController) GetContactUs()  {
		this.View("contactUs.tpl")
}

// GetPrivacy
// @router /app/privacy [get]
func (this *IndexController) GetPrivacy()  {
		this.View("privacy.tpl")
}

// GetAgreement
// @router /app/agreement [get]
func (this *IndexController) GetAgreement()  {
		this.View("agreement.tpl")
}