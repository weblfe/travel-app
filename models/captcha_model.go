package models

type CaptchaModel struct {
		BaseModel
}

func CaptchaModelOf() *CaptchaModel  {
		var model = new(CaptchaModel)
		model.Bind(model)
		model.Init()
		return model
}