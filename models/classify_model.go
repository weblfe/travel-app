package models

type ClassifyModel struct {
		BaseModel
}

func ClassifyModelOf() *ClassifyModel  {
		var model = new(ClassifyModel)
		model._Binder = model
		model.Init()
		return model
}