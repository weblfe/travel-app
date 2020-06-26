package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/libs"
		"time"
)

type UserModel struct {
		BaseModel
}

type User struct {
		Id                 bson.ObjectId `json:"id" bson:"_id"`
		UserNumId          int64         `json:"user_num_id" bson:"user_num_id"`
		UserName           string        `json:"username" bson:"username"`
		AvatarId           string        `json:"avatar_id,omitempty" bson:"avatar_id,omitempty"`
		NickName           string        `json:"nickname,omitempty" bson:"nickname,omitempty"`
		Password           string        `json:"password" bson:"password"`
		Mobile             string        `json:"mobile" bson:"mobile"`
		Email              string        `json:"email,omitempty" bson:"email,omitempty"`
		ResetPasswordTimes int           `json:"reset_password_times" bson:"reset_password_times"`
		RegisterWay        string        `json:"register_way" bson:"register_way"`
		AccessTokens       []string      `json:"access_tokens" bson:"access_tokens"`
		LastLoginAt        int64         `json:"last_login_at" bson:"last_login_at"`
		LastLoginLocation  string        `json:"last_login_location" bson:"last_login_location"`
		Status             int           `json:"status" bson:"status"`
		CreatedAt          time.Time     `json:"created_at" bson:"created_at"`
		UpdatedAt          time.Time     `json:"updated_at" bson:"updated_at"`
		DeletedAt          int64         `json:"deleted_at" bson:"deleted_at"`
}

func UserModelOf() *UserModel {
		var model = new(UserModel)
		model._Self = model
		model.Init()
		return model
}

func (this *User) Load(data map[string]interface{}) *User {
		for key, v := range data {
				this.Set(key, v)
		}
		return this
}

func (this *User) Set(key string, v interface{}) *User {
		switch key {
		case "user_num_id":
				this.UserNumId = v.(int64)
		case "username":
				this.UserName = v.(string)
		case "id":
				this.Id = v.(bson.ObjectId)
		case "password":
				if this.Password != "" {
						return this
				}
				if pass, ok := v.(string); ok {
						this.Password = libs.PasswordHash(pass)
				}
		case "nickname":
				this.NickName = v.(string)
		case "mobile":
				this.Mobile = v.(string)
		case "email":
				this.Email = v.(string)
		case "reset_password_times":
				this.ResetPasswordTimes = v.(int)
		case "status":
				this.Status = v.(int)
		case "access_tokens":
				if str, ok := v.(string); ok {
						this.AccessTokens = []string{str}
				}
				if str, ok := v.([]string); ok {
						this.AccessTokens = str
				}
		case "last_login_at":
				this.LastLoginAt = v.(int64)
		case "last_login_location":
				this.LastLoginLocation = v.(string)
		case "created_at":
				this.CreatedAt = v.(time.Time)
		case "updated_at":
				this.UpdatedAt = v.(time.Time)
		case "deleted_at":
				this.DeletedAt = v.(int64)
		}
		return this
}

func (this *User) Defaults() *User {
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.UserNumId == 0 {
				user := UserModelOf()
				this.UserNumId = libs.GetId(user.GetDatabaseName(), user.TableName(), user.GetConn())
		}
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now()
		}
		if this.Status == 0 {
				this.Status = 1
		}
		if this.UserName == "" && this.Mobile != "" {
				this.UserName = this.Mobile
		}
		if this.UserName == "" && this.Email != "" {
				this.UserName = this.Email
		}
		if this.Password == "" {
				this.Password = libs.PasswordHash(beego.AppConfig.DefaultString("default_password", "123456&Hex"))
		}
		return this
}

func (this *UserModel) CreateIndex() {
		// unique mobile
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"mobile"},
				Unique: true,
				Sparse: true,
		})
		// null unique email
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"email"},
				Unique: true,
				Sparse: true,
		})
		// null unique username
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"username"},
				Unique: true,
				Sparse: true,
		})
		_ = this.Collection().EnsureIndexKey("state")
		_ = this.Collection().EnsureIndexKey("nickname")
		_ = this.Collection().EnsureIndexKey("user_num_id")
		_ = this.Collection().EnsureIndexKey("avatar_id")
		_ = this.Collection().EnsureIndexKey("last_login_location", "last_login_at")
}

func (this *UserModel) TableName() string {
		return "users"
}
