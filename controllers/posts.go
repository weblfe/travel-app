package controllers

import "github.com/weblfe/travel-app/repositories"

type PostsController struct {
		BaseController
}

// 游记模块 controller
func PostsControllerOf() *PostsController  {
	 return new(PostsController)
}

// 发布游记
// @router /posts/create [post]
func (this *PostsController)Create()  {
		this.Send(repositories.NewPostsRepository(&this.BaseController.Controller).Create())
}

// 更新游记
// @router /posts/update [put]
func (this *PostsController)Update()  {
		this.Send(repositories.NewPostsRepository(&this.BaseController.Controller).Update())
}

// 列表我的
// @router /posts/lists/my [get]
func (this *PostsController)ListMy()  {
		this.Send(repositories.NewPostsRepository(&this.BaseController.Controller).Lists("my"))
}

// 更新
// @router /posts/lists/:address [get]
func (this *PostsController)ListByAddress()  {
		this.Send(repositories.NewPostsRepository(&this.BaseController.Controller).Lists("address"))
}

// 查询
// @router /posts/search  [get]
func (this *PostsController)Search()  {
		this.Send(repositories.NewPostsRepository(&this.BaseController.Controller).Lists("search"))
}

// 文章详情
// @router /posts/:id   [get]
func (this *PostsController)DetailById()  {
		this.Send(repositories.NewPostsRepository(&this.BaseController.Controller).GetById())
}

// 删除文章
// @router /posts/:id  [delete]
func (this *PostsController)RemoveById()  {
		this.Send(repositories.NewPostsRepository(&this.BaseController.Controller).RemoveId())
}
