package controllers

import (
		"github.com/weblfe/travel-app/repositories"
)

// 搜索
// @router /posts/search  [get]
func (this *PostsController) Search() {
		this.Send(repositories.NewPostsRepository(this).Lists("search"))
}

// 排行榜
// @router /posts/ranking  [get]
func (this *PostsController) Ranking() {
		this.Send(repositories.NewPostsRepository(this).GetRanking())
}

// 关注
// @router /posts/follows  [get]
func (this *PostsController) Follows() {
		this.Send(repositories.NewPostsRepository(this).GetFollows())
}

// 获取喜欢
// @router /posts/likes  [get]
func (this *PostsController) Likes() {
		this.Send(repositories.NewPostsRepository(this).GetLikes())
}

// 获取喜欢
// @router /posts/user/likes  [get]
func (this *PostsController) LikesQuery() {
	this.Send(repositories.NewPostsRepository(this).GetLikes(this.GetString("userId")))
}

// 获取推荐列表
// @router /posts/recommends  [get]
func (this *PostsController) Recommends() {
		this.Send(repositories.NewPostsRepository(this).ListsByPostType(this.GetString("type", "")))
}


// 获取喜欢
// @router /posts/all  [get]
func (this *PostsController) All() {
		this.Send(repositories.NewPostsRepository(this).GetAll())
}