package controllers

type PostsController struct {
		BaseController
}

// 验证码模块 controller
func PostsControllerOf() *PostsController  {
	 return new(PostsController)
}

// @router /posts/create [post]
func (this *PostsController)Create()  {

}

// @router /posts/update [put]
func (this *PostsController)Update()  {

}

// @router /posts/:id   [get]
func (this *PostsController)DetailById()  {

}

// @router /posts/:id  [delete]
func (this *PostsController)RemoveById()  {

}
