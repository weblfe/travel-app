package models

type PostsModel struct {
		BaseModel
}

func PostsModelOf() *PostsModel  {
		var model = new(PostsModel)
		model._Self = model
		model.Init()
		return model
}