package models

type ConfigModel struct {
		BaseModel
}

func ConfigModelOf() *ConfigModel  {
		var model = new(ConfigModel)
		model._Self = model
		model.Init()
		return model
}