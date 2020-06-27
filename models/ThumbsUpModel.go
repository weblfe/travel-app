package models

type ThumbsUpModel struct {
		BaseModel
}

func ThumbsUpModelOf() *ThumbsUpModel  {
		var model = new(ThumbsUpModel)
		model._Self = model
		model.Init()
		return model
}