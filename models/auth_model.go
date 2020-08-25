package models

type AuthModel struct {
		BaseModel
}

func AuthModelOf() *AuthModel  {
		var model = new(AuthModel)
		model._Binder = model
		model.Init()
		return model
}