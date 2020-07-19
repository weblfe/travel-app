package models

import (
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/transforms"
		"time"
)

type UserRoleTypeModel struct {
		BaseModel
}

type UserRoleType struct {
		Id        bson.ObjectId `json:"id" bson:"_id"`
		Name      string        `json:"name" bson:"name"`
		Role      int           `json:"number" json:"number"`
		CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
		dataClassImpl
}

func NewUserRoleType() *UserRoleType {
		var data = new(UserRoleType)
		data.Init()
		return data
}

func (this *UserRoleType) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *UserRoleType) data() bson.M {
		return bson.M{

		}
}

func (this *UserRoleType) setAttributes(data map[string]interface{}, safe ...bool) {

}

func (this *UserRoleType) setDefaults() {

}

func (this *UserRoleType) save() error {
		return nil
}
