package models

type RecommendModel struct {
		BaseModel
}

func RecommendModelOf() *RecommendModel  {
		var model = new(RecommendModel)
		model._Self = model
		model.Init()
		return model
}