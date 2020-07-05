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
		UserNumId          int64         `json:"userNumId" bson:"userNumId"`
		UserName           string        `json:"username" bson:"username"`
		AvatarId           string        `json:"avatarId,omitempty" bson:"avatarId,omitempty"`
		NickName           string        `json:"nickname,omitempty" bson:"nickname,omitempty"`
		PasswordHash       string        `json:"passwordHash" bson:"passwordHash"`
		Mobile             string        `json:"mobile" bson:"mobile"`
		Email              string        `json:"email,omitempty" bson:"email,omitempty"`
		ResetPasswordTimes int           `json:"resetPasswordTimes" bson:"resetPasswordTimes"`
		RegisterWay        string        `json:"registerWay" bson:"registerWay"`
		AccessTokens       []string      `json:"accessTokens" bson:"accessTokens"`
		LastLoginAt        int64         `json:"lastLoginAt" bson:"lastLoginAt"`
		LastLoginLocation  string        `json:"lastLoginLocation" bson:"lastLoginLocation"`
		Status             int           `json:"status" bson:"status"`
		Gender             int           `json:"gender" bson:"gender"`
		CreatedAt          time.Time     `json:"createdAt" bson:"createdAt"`
		UpdatedAt          time.Time     `json:"updatedAt" bson:"updatedAt"`
		DeletedAt          int64         `json:"deletedAt" bson:"deletedAt"`
}

const (
		UserTable        = "users"
		GenderUnknown    = 0 // 未知
		GenderMan        = 1 // 男
		GenderWoman      = 2 // 女
		GenderSecrecy    = 3 // 保密
		GenderBoth       = 4 // 中间人
		GenderSecrecyKey = "secrecy"
		GenderUnknownKey = "default"
		GenderManKey     = "man"
		GenderWomanKey   = "woman"
		GenderBothKey    = "both"
)

func UserModelOf() *UserModel {
		var model = new(UserModel)
		model._Self = model
		model.Init()
		return model
}

func NewUser() *User {
		var user = new(User)
		return user
}

func (this *User) Load(data map[string]interface{}) *User {
		for key, v := range data {
				this.Set(key, v)
		}
		return this
}

func (this *User) Set(key string, v interface{}) *User {
		switch key {
		case "userNumId":
				this.UserNumId = v.(int64)
		case "username":
				this.UserName = v.(string)
		case "id":
				this.Id = v.(bson.ObjectId)
		case "passwordHash":
				if this.PasswordHash != "" {
						return this
				}
				if pass, ok := v.(string); ok {
						this.PasswordHash = libs.PasswordHash(pass)
				}
		case "registerWay":
				this.RegisterWay = v.(string)
		case "nickname":
				this.NickName = v.(string)
		case "mobile":
				this.Mobile = v.(string)
		case "email":
				this.Email = v.(string)
		case "resetPasswordTimes":
				this.ResetPasswordTimes = v.(int)
		case "status":
				this.Status = v.(int)
		case "accessTokens":
				if str, ok := v.(string); ok {
						this.AccessTokens = []string{str}
				}
				if str, ok := v.([]string); ok {
						this.AccessTokens = str
				}
		case "lastLoginAt":
				this.LastLoginAt = v.(int64)
		case "lastLoginLocation":
				this.LastLoginLocation = v.(string)
		case "createdAt":
				this.CreatedAt = v.(time.Time)
		case "updatedAt":
				this.UpdatedAt = v.(time.Time)
		case "deletedAt":
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
		if this.NickName == "" && this.UserName != "" {
				this.NickName = this.UserName+"_nick"
		}
		if this.PasswordHash == "" {
				this.PasswordHash = libs.PasswordHash(beego.AppConfig.DefaultString("default_password", "123456&Hex"))
		}
		return this
}

func (this *User) M(filter ...func(m beego.M) beego.M) beego.M {
		data := beego.M{
				"id":                 this.Id.Hex(),
				"avatarId":           this.AvatarId,
				"gender":             this.Gender,
				"passwordHash":       this.PasswordHash,
				"username":           this.UserName,
				"nickname":           this.NickName,
				"registerWay":        this.RegisterWay,
				"mobile":             this.Mobile,
				"email":              this.Email,
				"userNumId":          this.UserNumId,
				"resetPasswordTimes": this.ResetPasswordTimes,
				"createdAt":          this.CreatedAt,
				"status":             this.Status,
				"lastLoginAt":        this.LastLoginAt,
				"lastLoginLocation":  this.LastLoginLocation,
				"deletedAt":          this.DeletedAt,
		}
		if len(filter) != 0 {
				for _, fn := range filter {
						data = fn(data)
				}
		}
		return data
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
		_ = this.Collection().EnsureIndexKey("gender")
		_ = this.Collection().EnsureIndexKey("nickname")
		_ = this.Collection().EnsureIndexKey("userNumId")
		_ = this.Collection().EnsureIndexKey("avatarId")
		_ = this.Collection().EnsureIndexKey("lastLoginLocation", "lastLoginAt")
}

func (this *UserModel) TableName() string {
		return UserTable
}

func GetGenderKey(gender int) string {
		switch gender {
		case GenderUnknown:
				return GenderUnknownKey
		case GenderMan:
				return GenderManKey
		case GenderWoman:
				return GenderWomanKey
		case GenderBoth:
				return GenderBothKey
		case GenderSecrecy:
				return GenderSecrecyKey
		}
		return GenderUnknownKey
}

func GetGenderEnum(gender string) int {
		switch gender {
		case GenderUnknownKey:
				return GenderUnknown
		case GenderManKey:
				return GenderMan
		case GenderWomanKey:
				return GenderWoman
		case GenderBothKey:
				return GenderBoth
		case GenderSecrecyKey:
				return GenderSecrecy
		}
		return GenderUnknown
}
