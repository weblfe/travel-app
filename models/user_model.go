package models

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
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
		Status             int           `json:"status" bson:"status"`                         // 用户状态 -1:拉黑,1:正常,2:禁用
		Gender             int           `json:"gender" bson:"gender"`                         // 用户性别 0:保密 1:男 2:女 3:😯
		Birthday           int64         `json:"birthday,omitempty" bson:"birthday,omitempty"` // 用户生日
		Address            string        `json:"address" bson:"address"`                       // 用户地址
		ThumbsUpTotal      int64         `json:"thumbsUpNum" bson:"thumbsUpNum"`               // 点赞总数
		InviteCode         string        `json:"inviteCode" bson:"inviteCode"`                 // 邀请码 6-64
		CreatedAt          time.Time     `json:"createdAt" bson:"createdAt"`                   // 创建时间 注册时间
		UpdatedAt          time.Time     `json:"updatedAt" bson:"updatedAt"`                   // 更新时间
		DeletedAt          int64         `json:"deletedAt" bson:"deletedAt"`                   // 删除时间戳
		dataClassImpl      `json:",omitempty" bson:",omitempty"`
}

// 头像信息
type AvatarInfo struct {
		AvatarUrl string `json:"avatarUrl"`
		AvatarId  string `json:"avatarId"`
}

// 地址信息
type AddressInfo struct {
		Address  string `json:"address"`
		Country  string `json:"country"`
		Province string `json:"province"`
		City     string `json:"city"`
		District string `json:"district"`
		Street   string `json:"street"`
}

// 公共用户相关信息
type UserProfile struct {
		Gender             int          `json:"gender"`             // 性别类型ID
		UserNumber         int64        `json:"userNumId"`          // 用户数字ID
		Intro              string       `json:"intro"`              // 简介
		BackgroundCoverUrl string       `json:"backgroundCoverUrl"` // 背景图片URL
		Avatar             *AvatarInfo  `json:"avatar"`             // 用户头像
		UserId             string       `json:"userId"`             // 用户唯一ID
		Address            *AddressInfo `json:"address"`            // 地址
		NickName           string       `json:"nickname"`           // 用户昵称
		GenderDesc         string       `json:"genderDesc"`         // 性别描述
		PostNumber         int64        `json:"postNum"`            // 用户作品数
		ThumbsUpNum        int64        `json:"thumbsUpNum"`        // 点赞数
		ThumbsUpNumTxt     string       `json:"thumbsUpNumTxt"`     // 点赞数字符串
		CommentNum         int64        `json:"commentNum"`         // 评论数
		CommentNumTxt      string       `json:"commentNumTxt"`      // 评论数字符串
		LikesNum           int64        `json:"likesNum"`           // 用户喜欢作品数量
		FollowNum          int64        `json:"followNum"`          // 用户关注数
		FansNum            int64        `json:"fansNum"`            // 用户粉丝数
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
		UserStatusOk     = 1  // 正常
		UserStatusForbid = 2  // 禁用
		UserStatusBack   = -1 // 拉黑
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
		model.Bind(model)
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
		case "thumbsUpNum":
				this.SetNumIntN(&this.ThumbsUpTotal, v)
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

func (this *User) IsForbid() bool {
		if this.Status == UserStatusBack || this.Status == UserStatusForbid {
				return true
		}
		return false
}

func (this *User) IsBlackList() bool {
		return this.Status == UserStatusBack
}

func (this *User) IsRootRole() bool {
		return this.Role == UserRootRole
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
				"thumbsUpNum":        this.ThumbsUpTotal,
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
		if err == nil {
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

func (this *UserModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *UserModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				// unique mobile
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"mobile"},
						Unique: true,
						Sparse: true,
				}))
				// null unique email
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"email"},
						Unique: true,
						Sparse: true,
				}))
				// null unique username
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"username"},
						Unique: true,
						Sparse: true,
				}))
				this.logs(doc.EnsureIndexKey("state"))
				this.logs(doc.EnsureIndexKey("gender"))
				this.logs(doc.EnsureIndexKey("address"))
				this.logs(doc.EnsureIndexKey("nickname"))
				this.logs(doc.EnsureIndexKey("userNumId"))
				this.logs(doc.EnsureIndexKey("avatarId"))
				this.logs(doc.EnsureIndexKey("lastLoginLocation", "lastLoginAt"))
		}
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

// 用户信息
func NewUserProfile(id string) *UserProfile {
		if id == "" {
				return nil
		}
		var (
				user    = NewUser()
				profile = new(UserProfile)
				err     = UserModelOf().GetById(id, user)
		)
		if err != nil {
				logs.Error(err)
				return nil
		}
		profile.Intro = user.Intro
		profile.Gender = user.Gender
		profile.UserId = user.Id.Hex()
		profile.NickName = user.NickName
		profile.UserNumber = user.UserNumId
		profile.GenderDesc = GenderText(user.Gender)
		profile.Avatar = NewAvatarInfo(user.AvatarId)
		profile.FansNum = GetUserFansNumber(profile.UserId)
		profile.LikesNum = GetUserLikeNumber(profile.UserId)
		profile.PostNumber = GetUserPostNumber(profile.UserId)
		profile.FollowNum = GetUserFollowNumber(profile.UserId)
		profile.CommentNum = GetUserCommentNumber(profile.UserId)
		profile.ThumbsUpNum = GetUserThumbsUpNumber(profile.UserId)
		profile.Address = NewAddressInfo(user.Address, profile.UserId)
		profile.CommentNumTxt = libs.BigNumberStringer(profile.CommentNum)
		profile.ThumbsUpNumTxt = libs.BigNumberStringer(profile.ThumbsUpNum)
		profile.BackgroundCoverUrl = GetUserBackgroundUrl(user.BackgroundCoverId, "mediaId")
		return profile
}

func NewAvatarInfo(avatarId string) *AvatarInfo {
		var avatar = new(AvatarInfo)
		avatar.AvatarId = avatarId
		var image = AttachmentModelOf().GetImageById(avatarId)
		if image != nil {
				avatar.AvatarUrl = image.Url
		}
		return avatar
}

func NewAddressInfo(address string, userId string) *AddressInfo {
		var addr = new(AddressInfo)
		addr.Country = "中国"
		if userId == "" {
				return addr
		}
		addr.Address = address
		addr.City = "广州"
		addr.Province = "广东"
		addr.District = ""
		addr.Street = ""
		return addr
}

// 作品评论数
func GetUserCommentNumber(userId string) int64 {
		var query = bson.M{
				"userId": userId,
				"status": beego.M{"$in": []int{StatusAuditPass, StatusWaitAudit}},
		}
		return int64(PostsModelOf().Sum(query, "commentNum"))
}

// 作品点赞数
func GetUserThumbsUpNumber(userId string) int64 {
		var query = bson.M{
				"userId": userId,
				"status": beego.M{"$in": []int{StatusAuditPass, StatusWaitAudit}},
		}
		return int64(PostsModelOf().Sum(query, "thumbsUpNum"))
}

// 获取用户背景封面图
func GetUserBackgroundUrl(id string, ty string) string {
		if id == "" {
				return ""
		}
		if ty == "mediaId" {
				var image = AttachmentModelOf().GetImageById(id)
				if image == nil {
						return ""
				}
				return image.Url
		}
		if ty == "userId" {
				var user = NewUser()
				if err := UserModelOf().GetById(id, user); err != nil {
						logs.Error(err)
						return ""
				}
				return GetUserBackgroundUrl(user.BackgroundCoverId, "mediaId")
		}
		return ""
}

// 作品可见作品数
func GetUserPostNumber(userId string) int64 {
		var query = bson.M{
				"userId": userId,
				"status": beego.M{"$in": []int{StatusAuditPass, StatusWaitAudit}},
		}
		return int64(PostsModelOf().Count(query))
}

// 该用户关注用户数
func GetUserFollowNumber(userId string) int64 {
		var query = beego.M{
				"status": StatusOk,
				"userId": userId,
		}
		return UserFocusModelOf().Count(query)
}

// 该用户喜欢作品用户数
func GetUserLikeNumber(userId string) int64 {
		var query = bson.M{
				"status": StatusOk,
				"type":   ThumbsTypePost,
				"userId": userId,
		}
		return int64(ThumbsUpModelOf().Count(query))
}

// 该用户粉丝数
func GetUserFansNumber(userId string) int64 {
		var query = beego.M{
				"status":      StatusOk,
				"focusUserId": userId,
		}
		return UserFocusModelOf().Count(query)
}

// 用户作品通过审核数
func GetUserPassPostNumber(userId string) int64 {
		var query = bson.M{
				"userId": userId,
				"status": StatusAuditPass,
		}
		return int64(PostsModelOf().Count(query))
}

// 用户作品未通过审核数
func GetUserNotPassPostNumber(userId string) int64 {
		var query = bson.M{
				"userId": userId,
				"status": StatusAuditUnPass,
		}
		return int64(PostsModelOf().Count(query))
}

// 获取用户作品下架数
func GetUserDownPostNumber(userId string) int64 {
		var query = bson.M{
				"userId": userId,
				"status": StatusAuditOff,
		}
		return int64(PostsModelOf().Count(query))
}

// 用户作品总数
func GetUserPostTotal(userId string) int64 {
		var query = bson.M{
				"userId": userId,
		}
		return int64(PostsModelOf().Count(query))
}
