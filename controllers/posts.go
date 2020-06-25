package controllers

import "github.com/astaxie/beego"

type PostsController struct {
		beego.Controller
}

// 验证码模块 controller
func PostsControllerOf() *PostsController  {
	 return new(PostsController)
}

// @route /posts/create [post]
func (this *PostsController)Create()  {

}

// @route /posts/update [put]
func (this *PostsController)Update()  {

}

// @route /posts/:id   [get]
func (this *PostsController)DetailById()  {

}

// @route /posts/:id  [delete]
func (this *PostsController)RemoveById()  {

}
