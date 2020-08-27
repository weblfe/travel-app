package models

type RecommendModel struct {
		BaseModel
}

func RecommendModelOf() *RecommendModel  {
		var model = new(RecommendModel)
		model.Bind(model)
		model.Init()
		return model
}