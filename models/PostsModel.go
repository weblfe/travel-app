package models

type PostsModel struct {
		BaseModel
}

func PostsModelOf() *PostsModel  {
		var model = new(PostsModel)
		model.Init()
		return model
}