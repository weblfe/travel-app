package controllers

import "github.com/weblfe/travel-app/repositories"

type ThumbsUpController struct {
		BaseController
}

func ThumbsUpControllerOf() *ThumbsUpController {
		var controller = new(ThumbsUpController)
		return controller
}

// @router /thumbs/up
func (this *ThumbsUpController) Post() {
		this.Send(repositories.NewThumbsUpRepository(this).Up())
}

// @router /thumbs/down
func (this *ThumbsUpController) Delete() {
		this.Send(repositories.NewThumbsUpRepository(this).Down())
}

// @router /thumbs/count
func (this *ThumbsUpController) Get() {
		this.Send(repositories.NewThumbsUpRepository(this).Count())
}
