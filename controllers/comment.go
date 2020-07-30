package controllers

import "github.com/weblfe/travel-app/repositories"

type CommentController struct {
		BaseController
}

func CommentControllerOf() *CommentController  {
		return new(CommentController)
}

// @router /comment/create [post]
func (this *CaptchaController)Create() {
   this.Send(repositories.NewCommentRepository(this).Create())
}