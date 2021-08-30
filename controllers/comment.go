package controllers

import "github.com/weblfe/travel-app/repositories"

type CommentController struct {
		BaseController
}

// CommentControllerOf 评论控制器
func CommentControllerOf() *CommentController  {
		return new(CommentController)
}

// Create
// @router /comment/create [post]
func (this *CaptchaController)Create() {
   this.Send(repositories.NewCommentRepository(this).Create())
}

// Lists
// @router /comment/list  [get]
func (this *CaptchaController)Lists() {
		this.Send(repositories.NewCommentRepository(this).Lists())
}
