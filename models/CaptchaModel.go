package models

type CaptchaModel struct {
		BaseModel
}

func CaptchaModelOf() *CaptchaModel  {
		var model = new(CaptchaModel)
		model._Self= model
		model.Init()
		return model
}