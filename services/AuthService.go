package services

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/cache"
		_ "github.com/astaxie/beego/cache/memcache"
		_ "github.com/astaxie/beego/cache/redis"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"time"
)

type AuthService interface {
		LoginByUserPassword(typ string, value string, password string, args ...interface{}) (*models.User, string, common.Errors)
		GetByAccessToken(string) (*models.User, common.Errors)
		Keep(token string, duration ...time.Duration)
		Token(user *models.User, args ...interface{}) string
		ReleaseByUserId(...string) bool
}

type AuthServiceImpl struct {
		BaseService
		storage   cache.Cache
		userModel *models.UserModel
}

const (
		IdKey                  = "id"
		CacheAtKey             = "cached_at"
		ExpiredAtKey           = "expired_at"
		AuthCacheDriverDefault = "redis"
		AuthAliveTime          = 7 * 24 * time.Hour
		AuthCacheDriverKey     = "auth_cache_driver"
		AuthCacheConfigKey     = "auth_cache_config"
		DispatchAccessToken    = "access_tokens"
		DispatchTokenKeep      = "keep"
		AuthCacheConfigDefault = `{"key":"access_token","conn":":6039","dbNum":"2","password":""}`
)

func AuthServiceOf() AuthService {
		var auth = new(AuthServiceImpl)
		auth.Init()
		return auth
}

func (this *AuthServiceImpl) Init() {
		this.init()
		this.userModel = models.UserModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return AuthServiceOf()
		}
		this.initStorage()
}

func (this *AuthServiceImpl) initStorage() {
		if this.storage != nil {
				return
		}
		driver := beego.AppConfig.DefaultString(AuthCacheDriverKey, AuthCacheDriverDefault)
		config := beego.AppConfig.DefaultString(AuthCacheConfigKey, AuthCacheConfigDefault)
		this.storage, _ = cache.NewCache(driver, config)
}

func (this *AuthServiceImpl) LoginByUserPassword(typ string, value string, password string, args ...interface{}) (*models.User, string, common.Errors) {
		var (
				user = &models.User{}
				err  = common.NewErrors()
		)
		errs := this.userModel.GetByKey(typ, value, user)
		if errs != nil {
				return nil, "", err.Set("msg", errs.Error()).Set("code", common.NotFound)
		}
		// 已被禁止用 || 已被删除
		if user.Status != 1 || user.DeletedAt != 0 {
				return nil, "", err.Set("msg", common.UserAccountForbid).Set("code", common.AccessForbid)
		}
		if !libs.PasswordVerify(user.PasswordHash, password) {
				return nil, "", err.Set("msg", common.PasswordOrAccountNotMatch).Set("code", common.VerifyNotMatch)
		}
		return user, this.Token(user, args...), nil

}

func (this *AuthServiceImpl) Token(user *models.User, args ...interface{}) string {
		user.LastLoginAt = time.Now().Unix()
		user.UpdatedAt = time.Now()
		token := libs.HashCode(*user)
		user.AccessTokens = append(user.AccessTokens, token)
		_ = this.userModel.Update(bson.M{"_id": user.Id}, user)
		data, alive := this.data(user)
		if data != nil {
				_ = this.GetCache().Put(token, data, alive)
		}
		// 异步更新 token 集合
		this.dispatch(user.Id.Hex(), DispatchAccessToken)
		return token
}

func (this *AuthServiceImpl) data(user *models.User) ([]byte, time.Duration) {
		now := time.Now().Unix()
		alive := this.getAliveTime()
		data := beego.M{IdKey: user.Id.Hex(), CacheAtKey: now, ExpiredAtKey: now + int64(alive)}
		if str, err := json.Marshal(&data); err == nil {
				return str, alive
		}
		return nil, alive
}

func (this *AuthServiceImpl) getAliveTime() time.Duration {
		return AuthAliveTime
}

func (this *AuthServiceImpl) dispatch(data string, name string) {
		switch name {
		case DispatchAccessToken:
				go this.updateUserAccessToken(data)
		case DispatchTokenKeep:
				go this.Keep(data)
		}
}

// 更新
func (this *AuthServiceImpl) updateUserAccessToken(uid string) {
		var (
				tokens  []string
				user    = &models.User{}
				storage = this.GetCache()
		)
		if err := this.userModel.GetById(uid, user); err != nil {
				arr := user.AccessTokens
				for _, token := range arr {
						if storage.IsExist(token) {
								tokens = append(tokens, token)
						}
				}
				user.AccessTokens = arr
				_ = this.userModel.Update(bson.M{"_id": user.Id}, user)
		}
}

// 获取用户数据 通过 token
func (this *AuthServiceImpl) GetByAccessToken(token string) (*models.User, common.Errors) {
		var (
				mapper beego.M
				user   = &models.User{}
		)
		mapper, ok := this.getTokenData(token)
		if !ok {
				return nil, common.NewErrors(common.InvalidTokenError, common.InvalidTokenCode)
		}
		id, ok := mapper["id"]
		if !ok {
				return nil, common.NewErrors(common.InvalidTokenError, common.InvalidTokenCode)
		}
		err := this.userModel.GetById(id.(string), user)
		if err == nil {
				return user, nil
		}
		return nil, common.NewErrors(err.Error(), common.NotFound)
}

// 获取token 数据
func (this *AuthServiceImpl) getTokenData(token string) (beego.M, bool) {
		var (
				mapper  beego.M
				storage = this.GetCache()
		)
		data := storage.Get(token)
		if data == nil {
				return nil, false
		}
		if d, ok := data.([]byte); ok {
				_ = json.Unmarshal(d, &mapper)
		}
		if len(mapper) == 0 {
				return nil, false
		}
		return mapper, true
}

// 保持登录token
func (this *AuthServiceImpl) Keep(token string, duration ...time.Duration) {
		if len(duration) == 0 {
				duration = append(duration, 24*time.Hour)
		}
		mapper, ok := this.getTokenData(token)
		if !ok {
				return
		}
		expiredAt, _ := mapper[ExpiredAtKey]
		cachedAt, _ := mapper[CacheAtKey]
		if expiredAt == nil {
				expiredAt = cachedAt.(int64) + int64(this.getAliveTime())
		}
		if expire, ok := expiredAt.(int64); ok {
				expire = int64(duration[0]) + expire
				mapper[ExpiredAtKey] = expire
				data, _ := json.Marshal(mapper)
				_ = this.GetCache().Put(token, data, time.Duration(expire-time.Now().Unix()))
		}
}

// 获取缓存
func (this *AuthServiceImpl) GetCache() cache.Cache {
		if this.storage != nil {
				return this.storage
		}
		this.initStorage()
		return this.storage
}

func (this *AuthServiceImpl) ReleaseByUserId(ids ...string) bool {
		if len(ids) == 0 {
				return false
		}
		var (
				userService = UserServiceOf()
				user        = userService.GetById(ids[0])
		)
		if user == nil {
				return false
		}
		if len(user.AccessTokens) == 0 {
				return false
		}
		for _, key := range user.AccessTokens {
				_ = this.GetCache().Delete(key)
		}
		user.AccessTokens = user.AccessTokens[0:0]
		data := beego.M{"accessTokens": user.AccessTokens, "modifies": []string{"accessTokens"}}
		_ = userService.UpdateByUid(user.Id.Hex(), data)
		return false
}
