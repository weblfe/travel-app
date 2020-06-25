package models

type CommentModel struct {
		BaseModel
}

func CommentModelOf() *CommentModel  {
		var model = new(CommentModel)
		model.Init()
		return model
}