package controllers

import "github.com/astaxie/beego"

type IndexController struct {
		beego.Controller
}

func IndexControllerOf() *IndexController {
		return new(IndexController)
}

// @route /  [get,post,delete,put]
func (this *IndexController)Index()  {
	 this.Data["json"] = beego.M{
	 		"code" : 0,
	 		"error" : "",
	 		"message" : "",
	 		"data" : beego.M{},
	 }
	 this.ServeJSON()
}
