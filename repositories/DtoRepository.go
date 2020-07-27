package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"time"
)

type DtoRepository struct{}

// 基础用户信息
type BaseUser struct {
		UserId     string  `json:"id"`       // 用户ID
		Nickname   string  `json:"nickname"` // 用户昵称
		AvatarInfo *Avatar `json:"avatar"`   // 用户头像
		Role       int     `json:"role"`     // 账号类型Id
		RoleDesc   string  `json:"roleDesc"` // 账号类型描述
}

// 简单用户信息
type SimpleUser struct {
		BaseUser
		InviteCode string `json:"inviteCode"` // 邀请码
		Intro      string `json:"intro"`      // 简介
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
		user.Role = data.Role
		user.RoleDesc = this.getRoleDesc(user.Role)
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
				user.AvatarInfo = this.GetAvatar(_user.AvatarId,_user.Gender)
				user.Nickname = _user.NickName
				user.InviteCode = _user.InviteCode
				user.Intro = _user.Intro
		case beego.M:

		case map[string]interface{}:

		}
		return user
}

func (this *DtoRepository) GetPrivacyUser(data interface{}) *SimpleUser {
		var user = new(SimpleUser)

		return user
}

func (this *DtoRepository) GetUser(data interface{}) *SimpleUser {
		var user = new(SimpleUser)

		return user
}
