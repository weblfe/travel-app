package models

type AttachmentModel struct {
		BaseModel
}

func AttachmentModelOf() *AttachmentModel  {
		var model = new(AttachmentModel)
		model.Init()
		return model
}