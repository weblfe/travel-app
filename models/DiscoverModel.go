package models

type DiscoverModel struct {
		BaseModel
}

func DiscoverModelOf() *DiscoverModel {
		var model = new(DiscoverModel)
		model._Binder = model
		model.Init()
		return model
}

func (this *DiscoverModel) Get()  {

}