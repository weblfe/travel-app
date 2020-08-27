package models

type DiscoverModel struct {
		BaseModel
}

func DiscoverModelOf() *DiscoverModel {
		var model = new(DiscoverModel)
		model.Bind(model)
		model.Init()
		return model
}

func (this *DiscoverModel) Get()  {

}