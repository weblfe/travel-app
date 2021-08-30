package controllers

import "github.com/weblfe/travel-app/repositories"

type TagsController struct {
		BaseController
}

func TagsControllerOf() *TagsController  {
		return new(TagsController)
}

// Lists
// @router  /tags  [get]
func (this *TagsController)Lists()  {
	 this.Send(repositories.NewTagRepository(this).GetPostTags())
}
