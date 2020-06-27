package models

type DiscoverModel struct {
		BaseModel
}

func DiscoverModelOf() *DiscoverModel {
		var model = new(DiscoverModel)
		model._Self = model
		model.Init()
		return model
}

func (this *DiscoverModel) Get()  {

}