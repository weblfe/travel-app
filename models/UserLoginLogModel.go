package models

import (
		"github.com/globalsign/mgo/bson"
		"time"
)

type UserLoginLogModel struct {
		BaseModel
}

type UserLoginLog struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`
		Uid           string        `json:"userId" bson:"userId"`
		LoginTime     time.Time     `json:"loginTime" bson:"loginTime"`
		Device        string        `json:"device" bson:"device"`
		Client        string        `json:"client" bson:"client"`
		LoginLocation string        `json:"loginLocation" bson:"loginLocation"`
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`
}

const (
		UserLoginLogTable = "user_login_log"
)

func UserLoginLogModelOf() *UserLoginLogModel {
		var model = new(UserLoginLogModel)
		model._Self = model
		model.Init()
		return model
}

func (this *UserLoginLogModel) TableName() string {
		return UserLoginLogTable
}
