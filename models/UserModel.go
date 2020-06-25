package models

type UserModel struct {
		BaseModel
}

func UserModelOf() *UserModel  {
		var model = new(UserModel)
		model.Init()
		return model
}