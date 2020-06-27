package models

type CommentModel struct {
		BaseModel
}

func CommentModelOf() *CommentModel  {
		var model = new(CommentModel)
		model._Self = model
		model.Init()
		return model
}