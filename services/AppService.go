package services

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/cache"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"strings"
)

type AppService interface {
		GetAppVersion(string) string
		GetAboutUs(...string) string
		GetPrivacy() string
		GetUserAgreement() string
		GetAppCustomers() []string
		GetAppInfos(...string) map[string]interface{}
}

type appServiceImpl struct {
		BaseService
		appModel    *models.AppModel
		configModel *models.ConfigModel
		store       cache.Cache
}

const (
		AppCache            = "app"
		AppCodeKey          = "appCode"
		ConfigAppScope      = "app"
		AppHelpUrlKey       = "appHelpUrl"
		AppAboutUsKey       = "appAboutUs"
		AppUserAgreementKey = "appUserAgreement"
		AppPrivacyKey       = "appPrivacy"
		AppCustomersKey     = "appCustomers"
		AppBuiltKey         = "appBuilt"
		AppAllegeEmail      = "allegeEmail"
)

func AppServiceOf() AppService {
		return newAppService()
}

// 创建新app服务
func newAppService() *appServiceImpl {
		var service = new(appServiceImpl)
		service.Init()
		return service
}

func (this *appServiceImpl) Init() {
		this.appModel = models.AppModelOf()
		this.configModel = models.ConfigModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return AppServiceOf()
		}
}

// 获取 App 版本号
func (this *appServiceImpl) GetAppVersion(device string) string {
		var (
				err error
				app = models.NewAppInfo()
		)
		err = this.appModel.NewQuery(bson.M{"device": device, "state": 1}).Sort("-publishTime").One(app)
		if err != nil {
				return "1.0.0"
		}
		return app.Version
}

// 获取关于我们
func (this *appServiceImpl) GetAboutUs(defaults ...string) string {
		var aboutUs = this.configModel.GetString(AppAboutUsKey)
		if aboutUs == "" {
				return defaults[0]
		}
		return aboutUs
}

// 获取隐私协议
func (this *appServiceImpl) GetPrivacy() string {
		return this.configModel.GetString(AppPrivacyKey)
}

// 获取用户协议
func (this *appServiceImpl) GetUserAgreement() string {
		return this.configModel.GetString(AppUserAgreementKey)
}

// 获取App 客户
func (this *appServiceImpl) GetAppCustomers() []string {
		var cnf = this.configModel.Get(AppCustomersKey, ConfigAppScope)
		if cnf.Value == nil {
				return []string{}
		}
		if str, ok := cnf.Value.(string); ok {
				return strings.SplitN(str, ",", -1)
		}
		if arr, ok := cnf.Value.([]string); ok {
				return arr
		}
		return []string{}
}

// 获取app 配置信息
func (this *appServiceImpl) GetAppInfos(driver ...string) map[string]interface{} {
		var (
				arr   []interface{}
				it    = this.GetConfig(ConfigAppScope)
				items = this.GetAppBuiltItems(driver...)
		)
		for _, it := range items {
				it.Init()
				arr = append(arr, it.M())
		}
		it[AppBuiltKey] = arr
		return it
}

// 获取构建版本信息内容
func (this *appServiceImpl) GetAppBuiltItems(drivers ...string) []*models.AppInfo {
		var (
				err   error
				count = len(drivers)
				items = make([]*models.AppInfo, 2)
				query = bson.M{"driver": bson.M{"$in": drivers}, "state": 1}
		)
		if count == 0 {
				count = 3
		}
		items = items[:0]
		err = this.appModel.NewQuery(query).Sort("-publishTime").Limit(count).All(&items)
		if err == nil {
				return items
		}
		return items
}

// 获取配置信息
func (this *appServiceImpl) GetConfig(typ ...string) beego.M {
		var (
				result = beego.M{}
				items  []models.Config
				arr    = make([]interface{}, 2)
		)
		if len(typ) == 0 {
				typ = append(typ, ConfigAppScope)
		}
		items = this.configModel.GetItemsByScope(typ[0])
		if len(items) == 0 {
				return result
		}
		arr = arr[:0]
		for _, it := range items {
				arr = append(arr, beego.M{
						"key":   it.Key,
						"value": it.Value,
						"title": it.Title,
				})
		}
		result["configs"] = arr
		return result
}

// 获取 App 三端 统一授权码
func (this *appServiceImpl) GetAppCode(defaults ...string) string {
		var code = this.configModel.GetString(AppCodeKey)
		if len(defaults) == 0 {
				defaults = append(defaults, "")
		}
		if code == "" {
				return defaults[0]
		}
		return code
}

// 获取帮助 Url
func (this *appServiceImpl) GetHelpUrl(defaults ...string) string {
		var url = this.configModel.GetString(AppHelpUrlKey)
		if len(defaults) == 0 {
				defaults = append(defaults, "")
		}
		if url == "" {
				return defaults[0]
		}
		return url
}

// 获取 申述 邮箱
func (this *appServiceImpl) GetAllegeEmail(defaults ...string) string {
		var email = this.configModel.GetString(AppAllegeEmail)
		if len(defaults) == 0 {
				defaults = append(defaults, "")
		}
		if email == "" {
				return defaults[0]
		}
		return email
}

func (this *appServiceImpl) SetAppCode(code string) bool {
		if err := this.configModel.Set(AppCodeKey, code); err == nil {
				return true
		}
		return false
}

// 获取缓存存储器
func (this *appServiceImpl) GetStore() cache.Cache {
		if this.store == nil {
				this.store = GetCacheService().Get(AppCache)
		}
		return this.store
}

func (this *appServiceImpl) code() string {
		return libs.RandomWord(6)
}

// 设置App Code
func (this *appServiceImpl) InitCode() {
		var code = this.GetAppCode()
		if code == "" {
				this.SetAppCode(this.code())
		}
}
