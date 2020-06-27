package models

type AttachmentModel struct {
		BaseModel
}

func AttachmentModelOf() *AttachmentModel  {
		var model = new(AttachmentModel)
		model._Self = model
		model.Init()
		return model
}