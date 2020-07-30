package repositories

import (
		"crypto/md5"
		"encoding/hex"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"sort"
		"strings"
		"time"
)

type DtoRepository struct {
		_Cache beego.M
}

// 基础用户信息
type BaseUser struct {
		UserId     string  `json:"id"`       // 用户ID
		Nickname   string  `json:"nickname"` // 用户昵称
		AvatarInfo *Avatar `json:"avatar"`   // 用户头像
}

// 简单用户信息
type SimpleUser struct {
		BaseUser
		InviteCode string `json:"inviteCode"` // 邀请码
		Intro      string `json:"intro"`      // 简介
		Role       int    `json:"role"`       // 账号类型Id
		RoleDesc   string `json:"roleDesc"`   // 账号类型描述
}

// 用户隐私数据
type PrivacyUser struct {
		SimpleUser
		Gender     int    `json:"gender"`     // 性别
		GenderDesc string `json:"genderDesc"` // 性别描述
		Birthday   int64  `json:"birthday"`   // 生日
		Address    string `json:"address"`    // 地址
}

// 用户数据
type User struct {
		PrivacyUser
		PasswordHash string    `json:"passwordHash"` // 密码
		Mobile       string    `json:"mobile"`       // 密码
		CreatedAt    time.Time `json:"createdAt"`    // 创建时间
		UpdatedAt    time.Time `json:"updatedAt"`    // 更新时间
}

type Avatar struct {
		Id        string `json:"id"`
		AvatarUrl string `json:"avatarUrl"`
}

var (
		_DTO = new(DtoRepository)
)

func GetDtoRepository() *DtoRepository {
		return _DTO
}

func (this *DtoRepository) GetUserById(id string) *BaseUser {
		var user = new(BaseUser)
		if id == "" {
				return user
		}
		var data = this.getUserService().GetById(id)
		if data == nil {
				return user
		}
		user.UserId = data.Id.Hex()
		user.Nickname = data.NickName
		user.AvatarInfo = this.GetAvatar(data.AvatarId, data.Gender)
		return user
}

func (this *DtoRepository) GetBaseUser(data *models.User) *BaseUser {
		var user = new(BaseUser)
		if data == nil {
				return user
		}
		user.UserId = data.Id.Hex()
		user.Nickname = data.NickName
		user.AvatarInfo = this.GetAvatar(data.AvatarId, data.Gender)
		return user
}

func (this *DtoRepository) GetBaseUserByMapper(data map[string]interface{}) *BaseUser {
		var user = new(BaseUser)
		if data == nil {
				return user
		}
		for key, v := range data {
				if str, ok := v.(string); ok && key == "id" {
						user.UserId = str
				}
				if id, ok := v.(bson.ObjectId); ok && key == "id" {
						user.UserId = id.Hex()
				}
				if str, ok := v.(string); ok && key == "nickname" {
						user.Nickname = str
				}
				if str, ok := v.(string); ok && key == "avatarId" {
						gender := data["gender"]
						if gender == nil {
								gender = 0
						}
						user.AvatarInfo = this.GetAvatar(str, gender.(int))
				}
		}
		return user
}

func (this *DtoRepository) GetAvatar(id string, gender int) *Avatar {
		var (
				avatar = new(Avatar)
				data   = this.getUserAvatarService().GetAvatarById(id, gender)
		)
		avatar.Id = data.Id
		avatar.AvatarUrl = data.Url
		return avatar
}

func (this *DtoRepository) getRoleDesc(role int) string {
		return this.GetUserRoleService().GetRoleDesc(role)
}

func (this *DtoRepository) GetUserRoleService() services.UserRoleService {
		return services.UserRoleServiceOf()
}

func (this *DtoRepository) getUserService() services.UserService {
		return services.UserServiceOf()
}

func (this *DtoRepository) getUserAvatarService() services.AvatarService {
		return services.AvatarServerOf()
}

func (this *DtoRepository) GetSimpleUserDetail(data interface{}) *SimpleUser {
		var user = new(SimpleUser)
		switch data.(type) {
		case *models.User:
				var _user = data.(*models.User)
				user.Role = _user.Role
				user.RoleDesc = this.getRoleDesc(user.Role)
				user.AvatarInfo = this.GetAvatar(_user.AvatarId, _user.Gender)
				user.Nickname = _user.NickName
				user.InviteCode = _user.InviteCode
				user.Intro = _user.Intro
				user.Role = _user.Role
				user.RoleDesc = this.getRoleDesc(user.Role)
		case beego.M:
				return this.GetUserByMapper(data.(beego.M))
		case map[string]interface{}:
				return this.GetUser(data.(map[string]interface{}))
		}
		return user
}

func (this *DtoRepository) GetSimpleUserDetailById(id string) *SimpleUser {
		var user = services.UserServiceOf().GetById(id)
		if user != nil {
				return this.GetSimpleUserDetail(user)
		}
		return nil
}

func (this *DtoRepository) GetPrivacyUser(data interface{}) *PrivacyUser {
		var user = new(PrivacyUser)
		switch data.(type) {
		case *models.User:
				var _user = data.(*models.User)
				user.RoleDesc = this.getRoleDesc(_user.Role)
				user.Role = _user.Role
				user.Intro = _user.Intro
				user.InviteCode = _user.InviteCode
				user.Nickname = _user.NickName
				user.Address = _user.Address
				user.AvatarInfo = this.GetAvatar(_user.AvatarId, _user.Gender)
				user.Birthday = _user.Birthday
				user.GenderDesc = models.GenderText(_user.Gender)
		case bson.M:
				var _user = data.(bson.M)
				return this.GetPrivacyUserByMapper(_user)
		case beego.M:
				var _user = data.(beego.M)
				return this.GetPrivacyUserByMapper(_user)
		case map[string]interface{}:
				var _user = data.(map[string]interface{})
				return this.GetPrivacyUserByMapper(_user)
		}
		return user
}

func (this *DtoRepository) GetPrivacyUserByMapper(data map[string]interface{}) *PrivacyUser {
		var user = new(PrivacyUser)
		for key, v := range data {
				if str, ok := v.(string); ok && key == "id" {
						user.UserId = str
				}
				if id, ok := v.(bson.ObjectId); ok && key == "id" {
						user.UserId = id.Hex()
				}
				if str, ok := v.(string); ok && key == "nickname" {
						user.Nickname = str
				}
				if str, ok := v.(string); ok && key == "inviteCode" {
						user.InviteCode = str
				}
				if str, ok := v.(string); ok && key == "intro" {
						user.Intro = str
				}
				if str, ok := v.(string); ok && key == "address" {
						user.Address = str
				}
				if gender, ok := v.(int); ok && key == "gender" {
						user.Gender = gender
				}
				if role, ok := v.(int); ok && key == "role" {
						user.Role = role
				}
				if roleDesc, ok := v.(string); ok && key == "roleDesc" {
						user.RoleDesc = roleDesc
				}
				if str, ok := v.(string); ok && key == "avatarId" {
						gender := data["gender"]
						if gender == nil {
								gender = 0
						}
						user.AvatarInfo = this.GetAvatar(str, gender.(int))
				}
		}
		if user.GenderDesc == "" {
				user.GenderDesc = models.GenderText(user.Gender)
		}
		return user
}

func (this *DtoRepository) GetUserByMapper(data beego.M) *SimpleUser {
		var user = new(SimpleUser)
		for key, v := range data {
				if str, ok := v.(string); ok && key == "id" {
						user.UserId = str
				}
				if id, ok := v.(bson.ObjectId); ok && key == "id" {
						user.UserId = id.Hex()
				}
				if str, ok := v.(string); ok && key == "nickname" {
						user.Nickname = str
				}
				if str, ok := v.(string); ok && key == "inviteCode" {
						user.InviteCode = str
				}
				if str, ok := v.(string); ok && key == "intro" {
						user.Intro = str
				}
				if role, ok := v.(int); ok && key == "role" {
						user.Role = role
				}
				if roleDesc, ok := v.(string); ok && key == "roleDesc" {
						user.RoleDesc = roleDesc
				}
				if str, ok := v.(string); ok && key == "avatarId" {
						gender := data["gender"]
						if gender == nil {
								gender = 0
						}
						user.AvatarInfo = this.GetAvatar(str, gender.(int))
				}
		}
		if user.Role != 0 && user.RoleDesc == "" {
				user.RoleDesc = this.getRoleDesc(user.Role)
		}
		return user
}

func (this *DtoRepository) GetUser(data map[string]interface{}) *SimpleUser {
		return this.GetUserByMapper(data)
}

func (this *DtoRepository) GetThumbsUpService() services.ThumbsUpService {
		return services.ThumbsUpServiceOf()
}

// 是否已点赞
func (this *DtoRepository) IsThumbsUp(postId string, userId string, status ...int) bool {
		if len(status) == 0 {
				status = append(status, 1)
		}
		var query = bson.M{
				"typeId": postId,
				"type":   "post",
				"userId": userId,
				"status": status[0],
		}
		return this.GetThumbsUpService().Exists(query)
}

func (this *DtoRepository) GC(key ...string) *DtoRepository {
		if len(key) == 0 {
				this._Cache = nil
				return this
		}
		for _, k := range key {
				delete(this._Cache, k)
		}
		return this
}

func (this *DtoRepository) Get(key string) interface{} {
		return this._Cache[key]
}

func (this *DtoRepository) Key(value ...interface{}) string {
		if len(value) == 0 {
				return ""
		}
		var (
				ins  = md5.New()
				keys = make([]string, 2)
		)
		keys = keys[:0]
		for _, v := range value {
				keys = append(keys, fmt.Sprintf("%v", v))
		}
		sort.Strings(keys)
		ins.Write([]byte(strings.Join(keys, "-")))
		return hex.EncodeToString(ins.Sum(nil))
}

func (this *DtoRepository) Cache(key string, v interface{}) *DtoRepository {
		if this._Cache == nil {
				this._Cache = beego.M{}
		}
		this._Cache[key] = v
		return this
}
