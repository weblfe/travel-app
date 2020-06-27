package models

type AuthModel struct {
		BaseModel
}

func AuthModelOf() *AuthModel  {
		var model = new(AuthModel)
		model._Self = model
		model.Init()
		return model
}