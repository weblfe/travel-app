package controllers

import "github.com/weblfe/travel-app/repositories"

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
