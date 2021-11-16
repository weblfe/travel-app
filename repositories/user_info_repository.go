package repositories

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"github.com/weblfe/travel-app/transports"
		"regexp"
		"time"
)

type UserInfoRepository interface {
		GetUserInfo() common.ResponseJson
		Search(string) common.ResponseJson
		ResetPassword() common.ResponseJson
		UpdateUserInfo() common.ResponseJson
		GetProfile(string) common.ResponseJson
		GetUserPublicInfo(string) common.ResponseJson
		GetUserFriends(ids ...string) common.ResponseJson
}

type UserInfoRepositoryImpl struct {
		userService         services.UserService
		ctx                 common.BaseRequestContext
		userBehaviorService services.UserBehaviorService
}

type UserNumbersObject struct {
		FollowNum   int64 `json:"followNum"`
		FansNum     int64 `json:"fansNum"`
		ThumbsUpNum int64 `json:"thumbsUpNum"`
}

type UserPublicInfo struct {
		UserNumbersObject
		BaseUser
		Intro              string `json:"intro"`              // 简介
		BackgroundCoverUrl string `json:"backgroundCoverUrl"` // 背景图
		Role               int    `json:"role"`               // 账号类型Id
		RoleDesc           string `json:"roleDesc"`           // 账号类型描述
		Gender             int    `json:"gender"`             // 性别
		GenderDesc         string `json:"genderDesc"`         // 性别描述
		Address            string `json:"address"`            // 地址
		IsFollowed         bool   `json:"isFollowed"`         // 是否已经关注
}

func NewUserInfoRepository(ctx common.BaseRequestContext) UserInfoRepository {
		var repository = new(UserInfoRepositoryImpl)
		repository.ctx = ctx
		repository.init()
		return repository
}

func (this *UserInfoRepositoryImpl) init() {
		this.userService = services.UserServiceOf()
		this.userBehaviorService = services.UserBehaviorServiceOf()
}

func (this *UserInfoRepositoryImpl) getDto() *DtoRepository {
		return GetDtoRepository()
}

// GetUserInfo 获取 全部用户信息 [个人]
func (this *UserInfoRepositoryImpl) GetUserInfo() common.ResponseJson {
		var (
				id string
				v  = this.ctx.Session(middlewares.AuthUserId)
		)
		if v != nil {
				id = v.(string)
		}
		if id == "" {
				return common.NewUnLoginResp(common.NewErrors(common.UnLoginCode, "请先登陆"))
		}
		user := this.userService.GetById(id)
		if user == nil || models.IsForbid(user) {
				return common.NewErrorResp(common.NewErrors(common.PermissionCode, "账号禁用状态"))
		}
		data := user.M(getBaseUserInfoTransform())

		return common.NewSuccessResp(beego.M{"user": data}, "获取成功")
}

func (this *UserInfoRepositoryImpl) ResetPassword() common.ResponseJson {
		var (
				ctx     = this.ctx.GetParent().Ctx
				request = new(transports.ResetPassword)
				userId  = this.ctx.GetParent().GetSession(middlewares.AuthUserId)
		)
		request.Load(ctx.Input)
		if request.Password == "" {
				return common.NewInvalidParametersResp("password参数缺失")
		}
		if err := this.PasswordCheck(request.Password); err != nil {
				return common.NewInvalidParametersResp(err)
		}
		// 登陆用户更新自己密码
		if request.Mobile == "" && request.Code == "" {
				if userId == "" || userId == nil {
						return common.NewInvalidParametersResp("异常请求")
				}
				request.UserId = userId.(string)
				if request.CurrentPassword == "" {
						return common.NewInvalidParametersResp("请输入当前登陆密码")
				}
				err := this.resetPasswordByCurrent(request.UserId, request.CurrentPassword, request.Password)
				if err == nil {
						return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix(), "tip": "请求重新登陆"}, "更新密码成功")
				}
				return common.NewFailedResp(common.ServiceFailed, "更新密码失败", err)
		}
		// 手机+验证码重置密码
		err := this.resetPasswordByMobileSms(request.Mobile, request.Code, request.Password)
		if err == nil {
				return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix(), "tip": "请求重新登陆"}, "更新密码成功")
		}
		return common.NewFailedResp(common.ServiceFailed, "更新密码失败", err)
}

// PasswordCheck 密码规则检查
func (this *UserInfoRepositoryImpl) PasswordCheck(pass string) common.Errors {
		if len(pass) < 6 {
				return common.NewErrors(common.InvalidParametersCode, "密码长度必须不少于6位")
		}
		var reg = regexp.MustCompile(`[\w@+.$&]{6,20}`)
		if !reg.MatchString(pass) {
				return common.NewErrors(common.InvalidParametersCode, "特殊越界,密码特殊字符是能包含(@,_,+,-,.,$,&)")
		}
		return nil
}

func (this *UserInfoRepositoryImpl) resetPasswordByCurrent(userId string, current string, password string) error {
		var (
				user = this.userService.GetById(userId)
		)
		if user == nil {
				return common.NewErrors("用户不存在")
		}
		if !libs.PasswordVerify(user.PasswordHash, current) {
				return common.NewErrors("密码不正确")
		}
		user.PasswordHash = ""
		// 更新密码
		user.Load(beego.M{
				"passwordHash": password,
		})
		// 仅更新密码
		data := beego.M{
				"passwordHash":       user.PasswordHash,
				"updatedAt":          time.Now(),
				"resetPasswordTimes": user.ResetPasswordTimes + 1,
				"modifies":           []string{"passwordHash", "updatedAt", "resetPasswordTimes"},
		}
		err := this.userService.UpdateByUid(userId, data)
		if err != nil {
				return common.NewErrors(err)
		}
		// 释放所有登陆token
		defer services.AuthServiceOf().ReleaseByUserId(userId)
		return nil
}

func (this *UserInfoRepositoryImpl) resetPasswordByMobileSms(mobile string, code string, password string) error {
		var (
				user = this.userService.GetByMobile(mobile)
		)
		if user == nil {
				return common.NewErrors("用户不存在")
		}
		// 校验
		if !services.SmsCodeServiceOf().Verify(mobile, code, "reset_password") {
				return common.NewErrors("验证码不正确")
		}
		user.PasswordHash = ""
		// 更新密码
		user.Load(beego.M{
				"passwordHash": password,
		})
		// 仅更新密码
		data := beego.M{
				"passwordHash":       user.PasswordHash,
				"updatedAt":          time.Now(),
				"resetPasswordTimes": user.ResetPasswordTimes + 1,
				"modifies":           []string{"passwordHash", "updatedAt", "resetPasswordTimes"},
		}
		err := this.userService.UpdateByUid(user.Id.Hex(), data)
		if err != nil {
				return common.NewErrors(err)
		}
		// 释放所有登陆token
		defer services.AuthServiceOf().ReleaseByUserId(user.Id.Hex())
		return nil
}

func (this *UserInfoRepositoryImpl) GetUserFriends(ids ...string) common.ResponseJson {
		if len(ids) == 0 {
				ids = append(ids, "")
		}
		var (
				users         = make([]*BaseUser, 2)
				currentUserId = getUserId(this.ctx)
				userId        = ids[0]
				page, count   = this.ctx.GetInt("page", 1), this.ctx.GetInt("count", 20)
				limit         = models.NewListParam(page, count)
		)
		if userId == "" {
				userId = currentUserId
		}
		if userId == "0" {
				return common.NewFailedResp(common.InvalidParametersCode, "用户ID缺失")
		}
		var userIds, meta = this.userBehaviorService.ListsByUserId(userId, limit)
		if userIds == nil || meta == nil {
				return common.NewFailedResp(common.RecordNotFound, "空")
		}
		users = users[:0]
		var dto = this.getDto()
		for _, user := range userIds {
				it := dto.GetUserById(user.FocusUserId.Hex())
				users = append(users, it)
		}
		if len(users) == 0 {
				return common.NewFailedResp(common.RecordNotFound, "空")
		}
		return common.NewSuccessResp(bson.M{"items": users, "meta": meta}, "获取成功")
}

func (this *UserInfoRepositoryImpl) UpdateUserInfo() common.ResponseJson {
		var (
				updateInfo    beego.M
				request       = this.ctx.GetParent()
				updateRequest = new(transports.UpdateUserRequest)
				userId        = request.GetSession(middlewares.AuthUserId)
				err           = updateRequest.Load(request.Ctx.Input.RequestBody)
		)
		if userId == nil || userId == "" {
				return common.NewUnLoginResp("登陆失效,请重新登陆！")
		}
		if err != nil {
				return common.NewErrorResp(common.NewErrors(err.Error(), common.ServiceFailed), "参数解析异常！")
		}
		if updateRequest.Empty() {
				return common.NewErrorResp(common.NewErrors(common.InvalidParametersError, common.InvalidParametersCode), "缺失请求参数")
		}
		updateInfo = updateRequest.M()
		if len(updateInfo) == 0 {
				return common.NewInvalidParametersResp("参数缺失!")
		}
		// 指定更新字段
		err = this.userService.UpdateByUid(userId.(string), updateInfo)
		if err == nil {
				return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix()}, "更新成功")
		}
		return common.NewFailedResp(common.ServiceFailed, err, "更新失败")
}

// GetUserNumbers 用户相关数值数据
func (this *UserInfoRepositoryImpl) GetUserNumbers(userId string) *UserNumbersObject {
		var info = new(UserNumbersObject)
		if !this.userService.Exists(beego.M{"_id": bson.ObjectIdHex(userId)}) {
				return info
		}
		info.FollowNum = this.userService.GetUserFollowCount(userId)
		info.FansNum = this.userService.GetUserFansCount(userId)
		info.ThumbsUpNum = this.userService.GetUserThumbsUpCount(userId)
		return info
}

// GetUserPublic 获取用户 公共信息
func (this *UserInfoRepositoryImpl) GetUserPublic(userId string) *UserPublicInfo {
		var info = new(UserPublicInfo)
		if !this.userService.Exists(beego.M{"_id": bson.ObjectIdHex(userId)}) {
				return info
		}
		var (
				dto  = GetDtoRepository()
				nums = this.GetUserNumbers(userId)
				user = this.userService.GetById(userId)
		)
		info.UserNumbersObject = *nums
		info.UserId = userId
		info.Intro = user.Intro
		info.Nickname = user.NickName
		info.Gender = user.Gender
		info.Address = user.Address
		info.Role = user.Role
		info.RoleDesc = user.GetRoleDesc(user.Role)
		info.GenderDesc = models.GenderText(user.Gender)
		info.AvatarInfo = dto.GetAvatar(user.AvatarId, user.Gender)
		info.IsFollowed = this.userBehaviorService.IsFollowed(getUserId(this.ctx), info.UserId)
		info.BackgroundCoverUrl = dto.GetUrlByAttachId(user.BackgroundCoverId)
		return info
}

func (this *UserInfoRepositoryImpl) GetUserPublicInfo(id string) common.ResponseJson {
		var user = this.GetUserPublic(id)
		if user == nil {
				return common.NewFailedResp(common.NotFound, "用户不存在")
		}
		return common.NewSuccessResp(bson.M{"user": user}, "获取成功")
}

func (this *UserInfoRepositoryImpl) Search(query string) common.ResponseJson {
		if query == "" {
				return common.NewFailedResp(common.InvalidParametersCode, "搜索参数异常")
		}
		var (
				ctx      = this.ctx.GetParent()
				page, _  = ctx.GetInt("page", 1)
				count, _ = ctx.GetInt("count", 20)
				limit    = models.NewListParam(page, count)
		)
		var items, meta = this.userService.SearchUserByNickName(query, limit)
		if items == nil || meta == nil {
				return common.NewFailedResp(common.NotFound, "空")
		}

		var (
				result     []beego.M
				transforms = getUserTransform(getUserId(this.ctx))
		)
		for _, it := range items {
				result = append(result, it.M(transforms, this.removes))
		}
		return common.NewSuccessResp(bson.M{"items": result, "meta": meta}, "获取成功")
}

func (this *UserInfoRepositoryImpl) removes(m beego.M) beego.M {
		var keys = []string{
				"resetPasswordTimes", "createdAt", "email",
				"birthday", "lastLoginAt", "inviteCode",
				"status", "lastLoginLocation",
				"username", "userNumId", "mobile",
				"backgroundCoverId",
		}
		for _, key := range keys {
				delete(m, key)
		}
		return m
}

func (this *UserInfoRepositoryImpl) GetProfile(userId string) common.ResponseJson {
		if userId == "" {
				userId = getUserId(this.ctx)
		}
		var data, err = this.userService.GetUserProfile(userId)
		if err == nil {
				return common.NewSuccessResp(beego.M{"user": data}, common.Success)
		}
		return common.NewFailedResp(common.NotFound, err.Error())
}
