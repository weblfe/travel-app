package models

type ClassifyModel struct {
		BaseModel
}

func ClassifyModelOf() *ClassifyModel  {
		var model = new(ClassifyModel)
		model.Bind(model)
		model.Init()
		return model
}