package repositories

import (
		"crypto/md5"
		"encoding/hex"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"sort"
		"strings"
		"sync"
		"time"
)

type DtoRepository struct {
		_Cache           beego.M
		_MaxCacheItemNum int
		_Table           cacheTable
		_Len             int
		_locker          sync.Mutex
		_Timer           *time.Ticker
		_Closer          chan byte
}

type cacheTable []*cache

// 缓存记录对象
type cache struct {
		Key      string
		CachedAt int64
		ExpireAt int64
		AccessAt int64
}

// BaseUser 基础用户信息
type BaseUser struct {
		UserId      string  `json:"userId"`      // 用户ID
		Nickname    string  `json:"nickname"`    // 用户昵称
		AvatarInfo  *Avatar `json:"avatar"`      // 用户头像
		UpdatedTime int64   `json:"updatedTime"` // 更新时间
}

// SimpleUser 简单用户信息
type SimpleUser struct {
		BaseUser
		InviteCode string `json:"inviteCode"` // 邀请码
		Intro      string `json:"intro"`      // 简介
		Role       int    `json:"role"`       // 账号类型Id
		RoleDesc   string `json:"roleDesc"`   // 账号类型描述
}

// PrivacyUser 用户隐私数据
type PrivacyUser struct {
		SimpleUser
		Gender     int    `json:"gender"`     // 性别
		GenderDesc string `json:"genderDesc"` // 性别描述
		Birthday   int64  `json:"birthday"`   // 生日
		Address    string `json:"address"`    // 地址
}

// User 用户数据
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
		_dtoSyncLock = sync.Once{}
		_DTO         *DtoRepository
)

const (
		DefaultMaxCacheItemNum    = 100
		DefaultCacheAliveDuration = 3 * time.Minute
)

func GetDtoRepository() *DtoRepository {
		if _DTO == nil {
				_dtoSyncLock.Do(func() {
						_DTO = newDto()
				})
		}
		return _DTO
}

func newDto() *DtoRepository {
		var dto = new(DtoRepository)
		dto._Cache = make(beego.M, 100)
		dto._Table = make([]*cache, 2)
		dto._Table = dto._Table[:0]
		dto._MaxCacheItemNum = DefaultMaxCacheItemNum
		dto._Timer = time.NewTicker(DefaultCacheAliveDuration)
		dto._Closer = make(chan byte, 2)
		dto.startGc()
		return dto
}

func (this *BaseUser) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = beego.M{
				"userId":   this.UserId,
				"nickname": this.Nickname,
				"avatar":   this.AvatarInfo,
		}
		if len(filters) == 0 {
				return data
		}
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}

// 修改最大缓存
func (this *DtoRepository) setMaxCache(max int) *DtoRepository {
		this._MaxCacheItemNum = max
		this._locker.Lock()
		defer this._locker.Unlock()
		this.check()
		return this
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
		user.UpdatedTime = data.UpdatedAt.Unix()
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
		user.UpdatedTime = data.UpdatedAt.Unix()
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
				if t, ok := v.(int64); ok && key == "updatedTime" {
						user.UpdatedTime = t
				}
				if t, ok := v.(time.Time); ok && key == "updatedTime" {
						user.UpdatedTime = t.Unix()
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

func (this *DtoRepository) GetUrlByAttachId(id string) string {
		var attach = services.AttachmentServiceOf().Get(id)
		return services.UrlTicketServiceOf().GetTicketUrlByAttach(attach)
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
				user.UpdatedTime = _user.UpdatedAt.Unix()
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
				user.UpdatedTime = _user.UpdatedAt.Unix()
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
				if t, ok := v.(int64); ok && key == "updatedTime" {
						user.UpdatedTime = t
				}
				if t, ok := v.(time.Time); ok && key == "updatedTime" {
						user.UpdatedTime = t.Unix()
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
				if t, ok := v.(int64); ok && key == "updatedTime" {
						user.UpdatedTime = t
				}
				if t, ok := v.(time.Time); ok && key == "updatedTime" {
						user.UpdatedTime = t.Unix()
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

func (this *DtoRepository) Stop() {
		this._Closer <- byte(1)
}

// IsThumbsUp 是否已点赞
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

// GC 回收数据
func (this *DtoRepository) GC(key ...string) *DtoRepository {
		this._locker.Lock()
		var total = len(this._Cache)
		if this._Len < total {
				this._Len = total
		}
		if len(key) == 0 {
				this._Len = 0
				this._Cache = beego.M{}
				this._Table = this._Table[:0]
				this._locker.Unlock()
				return this
		}
		for _, k := range key {
				delete(this._Cache, k)
				this._Table = this._Table.Delete(k)
				this._Len--
		}
		this._locker.Unlock()
		return this
}

// Get 获取缓存对象
func (this *DtoRepository) Get(key string) interface{} {
		this._locker.Lock()
		defer this._locker.Unlock()
		var v, ok = this._Cache[key]
		if ok {
				if res := this._Table.CheckAndUpdate(key); res == 0 {
						go this.GC(key)
						return nil
				}
		}
		return v
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

// Cache 缓存数据
func (this *DtoRepository) Cache(key string, v interface{}, alive ...time.Duration) *DtoRepository {
		this._locker.Lock()
		this.check()
		_, ok := this._Cache[key]
		if ok {
				this._Table.Update(key)
		}
		this._Cache[key] = v
		if len(alive) == 0 {
				alive = append(alive, DefaultCacheAliveDuration)
		}
		if !ok {
				this._Len++
				this._Table = this._Table.Append(key, int64(alive[0]))
		}
		this._locker.Unlock()
		return this
}

func (this *DtoRepository) Flash() *DtoRepository {
		this._locker.Lock()
		this.check()
		this._Len = 0
		this._Table.FlashALL()
		this._Cache = make(beego.M, 10)
		this._locker.Unlock()
		return this
}

// 获取cdn url
func (this *DtoRepository) getCdnUrl(url string, ty ...libs.OssUrlType) string {
		if len(ty) == 0 {
				ty = append(ty, libs.Row)
		}
		return libs.GetCdnUrl(url, ty[0])
}

// 检查是否触发 lru gc 条件
func (this *DtoRepository) check() {
		if (this._Len + this._MaxCacheItemNum/10) > this._MaxCacheItemNum {
				go this.lru()
		}
}

// 启动内部Gc 协程
func (this *DtoRepository) startGc() {
		go func() {
				var run = true
				for {
						select {
						case <-this._Timer.C:
								this.lru(-1)
								// 控制器
						case <-this._Closer:
								this._Timer.Stop()
								this.GC()
								run = false
								goto _DOC
						}
				_DOC:
						if !run {
								logs.Info("系统内部缓存GC,stop...")
								break
						}
				}
		}()
}

// Len 缓存数量
func (this *DtoRepository) Len() int {
		return this._Len
}

// 最近最少使用到的
func (this *DtoRepository) lru(num ...int) {
		// 无缓存
		if this._Len == 0 && this._Table.Len() == 0 {
				return
		}
		var arr = this._Table[:]
		sort.Sort(&arr)
		if len(num) == 0 {
				num = append(num, arr.Len()/2)
		}

		// 定时Gc
		if num[0] == -1 {
				this.expireGc(arr)
				return
		}
		// 手动回收对应空间Gc
		this.gcBySize(arr, num[0])
}

// 获取自动收回
func (this *DtoRepository) expireGc(table cacheTable) {
		var (
				deleteKeys []string
				now        = time.Now().Unix()
		)
		for _, cache := range table {
				if cache.ExpireAt <= now && cache.ExpireAt != -1 {
						deleteKeys = append(deleteKeys, cache.Key)
				}
		}
		if len(deleteKeys) > 0 {
				logs.Info(fmt.Sprintf("回收数据：%v", deleteKeys))
				this.GC(deleteKeys...)
		}
}

// 添加关注状态
func (this *DtoRepository) appendFollowStatus(userId string) func(m beego.M) beego.M {
		return func(m beego.M) beego.M {
				m["isFollowed"] = false
				if id, ok := m["userId"]; ok {
						m["isFollowed"] = this.IsFollowed(userId, id.(string))
				}
				return m
		}
}

// IsFollowed 是否已关注
func (this *DtoRepository) IsFollowed(userId, followerUserId string) bool {
		return services.UserBehaviorServiceOf().IsFollowed(userId, followerUserId)
}

// 手动回收
func (this *DtoRepository) gcBySize(table cacheTable, size int) {
		var total = table.Len()
		if size > total {
				size = total
		}
		logs.Info("系统内部缓存GC,checking...")
		this._Table = table[:size]
		this._Len = this._Table.Len()
}

func (this cacheTable) Len() int {
		return len(this)
}

func (this cacheTable) Append(key string, alive ...int64) cacheTable {
		if len(alive) == 0 {
				alive = append(alive, int64(DefaultCacheAliveDuration))
		}
		var (
				now    = time.Now().Unix()
				expire = alive[0]
				it     = cache{
						Key:      key,
						CachedAt: now,
						ExpireAt: now + expire,
						AccessAt: now,
				}
		)
		return append(this, &it)
}

func (this cacheTable) Update(key string) {
		for _, it := range this {
				if it.Key == key {
						it.AccessAt = time.Now().Unix()
				}
		}
}

func (this cacheTable) Delete(key string) cacheTable {
		for i, it := range this {
				if it.Key == key {
						return append(this[:i], this[i+1:]...)
				}
		}
		return this
}

func (this cacheTable) Less(i, j int) bool {
		var ita, itb = this[i], this[j]
		if ita.AccessAt-ita.CachedAt < itb.AccessAt-itb.CachedAt {
				if ita.CachedAt < itb.CachedAt {
						return true
				}
		}
		return false
}

func (this cacheTable) Swap(i, j int) {
		this[i], this[j] = this[j], this[i]
}

func (this cacheTable) FlashALL() cacheTable {
		return this[:0]
}

// CheckAndUpdate 检查和更新访问
func (this cacheTable) CheckAndUpdate(key string) int {
		var now = time.Now().Unix()
		for _, it := range this {
				if it.ExpireAt <= now {
						return 0
				}
				if it.Key == key {
						it.AccessAt = time.Now().Unix()
						return 1
				}
		}
		return -1
}
