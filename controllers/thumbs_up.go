package controllers

import "github.com/weblfe/travel-app/repositories"

type ThumbsUpController struct {
		BaseController
}

// ThumbsUpControllerOf 点赞控控制器
func ThumbsUpControllerOf() *ThumbsUpController {
		var controller = new(ThumbsUpController)
		return controller
}

// Post
// @router /thumbsUp  [post]
func (this *ThumbsUpController) Post() {
		this.Send(repositories.NewThumbsUpRepository(this).Up())
}

// Delete
// @router /thumbsUp [delete]
func (this *ThumbsUpController) Delete() {
		this.Send(repositories.NewThumbsUpRepository(this).Down())
}

// Get
// @router /thumbsUp/count [get]
func (this *ThumbsUpController) Get() {
		this.Send(repositories.NewThumbsUpRepository(this).Count())
}
