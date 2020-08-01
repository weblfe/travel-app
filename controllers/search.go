package controllers

import "github.com/weblfe/travel-app/repositories"

// 搜索
// @router /posts/search  [get]
func (this *PostsController) Search() {
		this.Send(repositories.NewPostsRepository(this).Lists("search"))
}