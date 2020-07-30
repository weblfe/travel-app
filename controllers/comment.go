package controllers

import "github.com/weblfe/travel-app/repositories"

type CommentController struct {
		BaseController
}

// @router /comment/create [post]
func (this *CaptchaController)Create() {
   this.Send(repositories.NewCommentRepository(this).Create())
}