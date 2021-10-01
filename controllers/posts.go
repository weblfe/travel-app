package controllers

import (
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/repositories"
)

type PostsController struct {
		BaseController
}

// PostsControllerOf 游记模块 controller
func PostsControllerOf() *PostsController {
		return new(PostsController)
}

// Create 发布游记
// @router /posts/create [post]
func (this *PostsController) Create() {
		this.Send(repositories.NewPostsRepository(this).Create())
}

// Update 更新游记
// @router /posts/update [put]
func (this *PostsController) Update() {
		this.Send(repositories.NewPostsRepository(this).Update())
}

// ListMy 列表我的
// @router /posts/lists/my [get]
func (this *PostsController) ListMy() {
		this.Send(repositories.NewPostsRepository(this).Lists("my"))
}

// ListBy 罗列作品信息列表 by tags
// @router /posts/lists [get]
func (this *PostsController) ListBy() {
		var typ = this.GetString("type", "tags")
		this.Send(repositories.NewPostsRepository(this).Lists(typ))
}

// ListUserPosts 其他用户作品
// @router /posts/users/:userId [get]
func (this *PostsController) ListUserPosts() {
		var data, ok = this.GetParam(":userId")
		// 不正常的用户ID
		if ok && data != nil && data != "" {
				this.Send(repositories.NewPostsRepository(this).Lists("user", data.(string)))
				return
		}
		this.Send(common.NewErrorResp(common.NewErrors(common.InvalidParametersCode, "异常ID"), common.InvalidParametersError))
}

// ListByAddress 通过地址罗列
// @router /posts/address/:address [get]
func (this *PostsController) ListByAddress() {
		this.Send(repositories.NewPostsRepository(this).Lists("address"))
}

// DetailById 文章详情
// @router /posts/:id   [get]
func (this *PostsController) DetailById() {
		this.Send(repositories.NewPostsRepository(this).GetById())
}

// RemoveById 删除文章
// @router /posts/:id  [delete]
func (this *PostsController) RemoveById() {
		this.Send(repositories.NewPostsRepository(this).RemoveId())
}

// Audit 审核文章
// @router /posts/audit  [post]
func (this *PostsController) Audit() {
		this.Send(repositories.NewPostsRepository(this).Audit())
}

// AutoCover 自动截图
// @router /posts/video/cover  [post]
func (this *PostsController) AutoCover() {
		this.Send(repositories.NewPostsRepository(this).AutoVideosCover())
}
