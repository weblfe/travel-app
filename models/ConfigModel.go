package models

type ConfigModel struct {
		BaseModel
}

func ConfigModelOf() *ConfigModel  {
		var model = new(ConfigModel)
		model.Init()
		return model
}