package models

import (
		"fmt"
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
		Id                 bson.ObjectId `json:"id" bson:"_id"`                                // 唯一ID
		UserNumId          int64         `json:"userNumId" bson:"userNumId"`                   // 用户注册序号
		Role               int           `json:"role" bson:"role"`                             // 用户类型
		UserName           string        `json:"username" bson:"username"`                     // 用户名唯一
		Intro              string        `json:"intro" bson:"intro"`                           // 个人简介
		BackgroundCoverId  string        `json:"backgroundCoverId" bson:"backgroundCoverId"`   // 个人也背景
		AvatarId           string        `json:"avatarId,omitempty" bson:"avatarId,omitempty"` // 头像ID
		NickName           string        `json:"nickname,omitempty" bson:"nickname,omitempty"` // 昵称
		PasswordHash       string        `json:"passwordHash" bson:"passwordHash"`             // 密码密码
		Mobile             string        `json:"mobile,omitempty" bson:"mobile,omitempty"`     // 手机号
		Email              string        `json:"email,omitempty" bson:"email,omitempty"`       // 邮箱
		ResetPasswordTimes int           `json:"resetPasswordTimes" bson:"resetPasswordTimes"` // 重置密码次数
		RegisterWay        string        `json:"registerWay" bson:"registerWay"`               // 注册方式
		AccessTokens       []string      `json:"accessTokens" bson:"accessTokens"`             // 授权临牌集合
		LastLoginAt        int64         `json:"lastLoginAt" bson:"lastLoginAt"`               // 最近一次登陆时间
		LastLoginLocation  string        `json:"lastLoginLocation" bson:"lastLoginLocation"`   // 最近一次登陆定位
		Status             int           `json:"status" bson:"status"`                         // 用户状态 1:正常
		Gender             int           `json:"gender" bson:"gender"`                         // 用户性别 0:保密 1:男 2:女 3:😯
		Birthday           int64         `json:"birthday,omitempty" bson:"birthday,omitempty"` // 用户生日
		Address            string        `json:"address" bson:"address"`                       // 用户地址
		InviteCode         string        `json:"inviteCode" bson:"inviteCode"`                 // 邀请码 6-64
		CreatedAt          time.Time     `json:"createdAt" bson:"createdAt"`                   // 创建时间 注册时间
		UpdatedAt          time.Time     `json:"updatedAt" bson:"updatedAt"`                   // 更新时间
		DeletedAt          int64         `json:"deletedAt" bson:"deletedAt"`                   // 删除时间戳
		dataClassImpl      `json:",omitempty" bson:",omitempty"`
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

var (
		genderMapper = map[int]string{
				GenderUnknown: "未知",
				GenderMan:     "男",
				GenderWoman:   "女",
				GenderSecrecy: "保密",
				GenderBoth:    "中间人",
		}
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

func GenderText(gender int) string {
		return genderMapper[gender]
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
				this.SetNumIntN(&this.UserNumId, v)
		case "username":
				this.SetString(&this.UserName, v)
		case "intro":
				fallthrough
		case "Intro":
				this.SetString(&this.Intro, v)
		case "id":
				this.SetObjectId(&this.Id, v)
		case "passwordHash":
				if this.PasswordHash != "" {
						return this
				}
				if pass, ok := v.(string); ok {
						this.PasswordHash = libs.PasswordHash(pass)
				}
		case "registerWay":
				this.SetString(&this.RegisterWay, v)
		case "nickname":
				this.SetString(&this.NickName, v)
		case "mobile":
				this.SetString(&this.Mobile, v)
		case "email":
				this.SetString(&this.Email, v)
		case "resetPasswordTimes":
				this.SetNumInt(&this.ResetPasswordTimes, v)
		case "status":
				this.SetNumInt(&this.Status, v)
		case "accessTokens":
				if str, ok := v.(string); ok {
						this.AccessTokens = []string{str}
				}
				if str, ok := v.([]string); ok {
						this.AccessTokens = str
				}
		case "lastLoginAt":
				this.SetNumIntN(&this.LastLoginAt, v)
		case "role":
				this.SetNumInt(&this.Role, v)
		case "lastLoginLocation":
				this.SetString(&this.LastLoginLocation, v)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
		case "inviteCode":
				this.SetString(&this.InviteCode, v)
		case "updatedAt":
				this.SetTime(&this.UpdatedAt, v)
		case "deletedAt":
				this.SetNumIntN(&this.DeletedAt, v)
		}
		return this
}

func (this *User) Defaults() *User {
		if this.Id == "" {
				this.Id = this.GetId()
		}
		if this.UserNumId == 0 {
				this.UserNumId = this.GetUserNumId()
		}
		if this.CreatedAt.IsZero() {
				this.CreatedAt = this.GetNow()
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
		if this.Mobile == "" && this.UserName != "" {
				this.Mobile = this.GetMobile()
		}
		if this.NickName == "" && this.UserName != "" {
				this.NickName = this.GetNickName()
		}
		if this.PasswordHash == "" {
				this.PasswordHash = this.GetPasswordHash()
		}
		if this.InviteCode == "" {
				this.InviteCode = this.GetInviteCode()
		}
		return this
}

func (this *User) GetPasswordHash() string {
		if this.PasswordHash == "" {
				return libs.PasswordHash(beego.AppConfig.DefaultString("default_password", "123456&Hex"))
		}
		return this.PasswordHash
}

func (this *User) GetMobile() string {
		if this.Mobile != "" {
				return this.Mobile
		}
		if libs.IsCnMobile(this.UserName) || libs.IsMobile(this.UserName) {
				return this.UserName
		}
		return ""
}

func (this *User) GetUserNumId() int64 {
		if this.UserNumId != 0 {
				return this.UserNumId
		}
		user := UserModelOf()
		return libs.GetId(user.GetDatabaseName(), user.TableName(), user.GetConn())
}

func (this *User) GetNickName() string {
		if this.NickName != "" {
				return this.NickName
		}
		return this.UserName + "_nick"
}

func (this *User) GetInviteCode(refresh ...bool) string {
		if len(refresh) == 0 {
				refresh = append(refresh, false)
		}
		if refresh[0] {
				return libs.Md5(fmt.Sprintf("%d", time.Now().Unix()))
		}
		if this.InviteCode == "" {
				return libs.Md5(fmt.Sprintf("%d", time.Now().Unix()))
		}
		return this.InviteCode
}

func (this *User) GetAddress(typ ...int) string {
		var addr = NewUserAddress()
		if this.Address != "" {
				return this.Address
		}
		if len(typ) == 0 {
				typ = append(typ, AddressTypeRegister)
		}
		_ = UserAddressModelOf().FindOne(bson.M{"userId": this.Id.Hex(), "type": typ[0]}, addr)
		return addr.String()
}

func (this *User) M(filter ...func(m beego.M) beego.M) beego.M {
		data := beego.M{
				"id":                 this.Id.Hex(),
				"avatarId":           this.AvatarId,
				"gender":             this.Gender,
				"role":               this.Role,
				"roleDesc":           this.GetRoleDesc(this.Role),
				"genderDesc":         GenderText(this.Gender),
				"passwordHash":       this.PasswordHash,
				"username":           this.UserName,
				"nickname":           this.NickName,
				"registerWay":        this.RegisterWay,
				"mobile":             this.Mobile,
				"email":              this.Email,
				"intro":              this.Intro,
				"backgroundCoverId":  this.BackgroundCoverId,
				"userNumId":          this.UserNumId,
				"resetPasswordTimes": this.ResetPasswordTimes,
				"status":             this.Status,
				"lastLoginAt":        this.LastLoginAt,
				"birthday":           this.Birthday,
				"createdAt":          this.CreatedAt.Unix(),
				"address":            this.GetAddress(),
				"inviteCode":         this.InviteCode,
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

func (this *User) Save() error {
		var (
				id    = this.Id.Hex()
				tmp   = new(User)
				model = UserModelOf()
				err   = model.GetById(id, tmp)
		)
		if err != nil {
				return model.UpdateById(id, this.M(func(m beego.M) beego.M {
						delete(m, "id")
						delete(m, "createdAt")
						m["updatedAt"] = time.Now().Local()
						return m
				}))
		}
		return model.Add(this.Defaults())
}

// 获取角色描述
func (this *User) GetRoleDesc(role int) string {
		return UserRolesConfigModelOf().GetRoleName(role)
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
		_ = this.Collection().EnsureIndexKey("address")
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

func IsForbid(data *User) bool {
		return data.DeletedAt != 0 || data.Status != 1
}
