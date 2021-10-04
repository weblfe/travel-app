package controllers

import (
		"github.com/weblfe/travel-app/repositories"
)

// Search 搜索
// @router /posts/search  [get]
func (this *PostsController) Search() {
		this.Send(repositories.NewPostsRepository(this).Lists("search"))
}

// Ranking 排行榜
// @router /posts/ranking  [get]
func (this *PostsController) Ranking() {
		this.Send(repositories.NewPostsRepository(this).GetRanking())
}

// Follows 关注
// @router /posts/follows  [get]
// @router /posts/follow  [get]
func (this *PostsController) Follows() {
		this.Send(repositories.NewPostsRepository(this).GetFollows())
}

// Likes 获取喜欢
// @router /posts/likes  [get]
func (this *PostsController) Likes() {
		this.Send(repositories.NewPostsRepository(this).GetLikes())
}

// LikesQuery 获取喜欢
// @router /posts/user/likes  [get]
func (this *PostsController) LikesQuery() {
	this.Send(repositories.NewPostsRepository(this).GetLikes(this.GetString("userId")))
}

// Recommends 获取推荐列表
// @router /posts/recommends  [get]
func (this *PostsController) Recommends() {
		this.Send(repositories.NewPostsRepository(this).ListsByPostType(this.GetString("type", "")))
}

// All 获取喜欢
// @router /posts/all  [get]
func (this *PostsController) All() {
		this.Send(repositories.NewPostsRepository(this).GetAll())
}