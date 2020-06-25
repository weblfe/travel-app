package models

type AppModel struct {
		BaseModel
}

func AppModelOf() *AppModel {
		var model = new(AppModel)
		model.Init()
		return model
}

