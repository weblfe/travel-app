package models

type AuthModel struct {
		BaseModel
}

func AuthModelOf() *AuthModel  {
		var model = new(AuthModel)
		model.Bind(model)
		model.Init()
		return model
}