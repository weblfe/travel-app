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
		Uid           string        `json:"user_id" bson:"user_id"`
		LoginTime     time.Time     `json:"login_time" bson:"login_time"`
		Device        string        `json:"device" bson:"device"`
		Client        string        `json:"client" bson:"client"`
		LoginLocation string        `json:"login_location" bson:"login_location"`
		CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
}

func (this *UserLoginLogModel) TableName() string {
		return "user_login_log"
}
