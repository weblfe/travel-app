package models

type AppModel struct {
		BaseModel
}

func AppModelOf() *AppModel {
		var model = new(AppModel)
		model._Self = model
		model.Init()
		return model
}

