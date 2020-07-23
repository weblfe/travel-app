package controllers

import "github.com/weblfe/travel-app/repositories"

type ThumbsUpController struct {
		BaseController
}

func ThumbsUpControllerOf() *ThumbsUpController {
		var controller = new(ThumbsUpController)
		return controller
}

// @router /thumbsUp  [post]
func (this *ThumbsUpController) Post() {
		this.Send(repositories.NewThumbsUpRepository(this).Up())
}

// @router /thumbsUp [delete]
func (this *ThumbsUpController) Delete() {
		this.Send(repositories.NewThumbsUpRepository(this).Down())
}

// @router /thumbsUp/count [get]
func (this *ThumbsUpController) Get() {
		this.Send(repositories.NewThumbsUpRepository(this).Count())
}
