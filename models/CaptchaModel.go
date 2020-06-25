package models

type CaptchaModel struct {
		BaseModel
}

func CaptchaModelOf() *CaptchaModel  {
		var model = new(CaptchaModel)
		model.Init()
		return model
}