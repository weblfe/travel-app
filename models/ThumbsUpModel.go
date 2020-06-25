package models

type ThumbsUpModel struct {
		BaseModel
}

func ThumbsUpModelOf() *ThumbsUpModel  {
		var model = new(ThumbsUpModel)
		model.Init()
		return model
}