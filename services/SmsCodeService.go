package services

import (
		"fmt"
		"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/cache"
		_ "github.com/astaxie/beego/cache/memcache"
		_ "github.com/astaxie/beego/cache/redis"
		"github.com/astaxie/beego/logs"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"os"
		"strconv"
		"strings"
		"sync"
		"time"
)

type SmsCodeService interface {
		Send(mobile string, content string, extras map[string]string) error
		SendCode(mobile string, typ string, extras map[string]string) (string, error)
		Verify(mobile string, code string, typ string) bool
		Storage(mobile string, data string, typ string, timeout time.Duration)
}

type smsServiceMocker interface {
		SmsCodeService
		SetContext(ctx SmsCodeService)
}

type SmsCodeServiceAliCloudImpl struct {
		client  *dysmsapi.Client
		storage cache.Cache
		mock    smsServiceMocker
		BaseService
}

const (
		DySmsApiRegionId        = "dysms_api_region_id"
		DySmsAccessKeyId        = "dysms_access_key_id"
		DySmsAccessKeySecret    = "dysms_access_key_secret"
		DefaultDySmsApiRegionId = "cn-hangzhou"
		SmsCacheDriverKey       = "sms_cache_driver"
		SmsCacheDriverDefault   = "redis"
		SmsCacheConfigKey       = "sms_cache_config"
		SmsCacheConfigDefault   = `{"key":"sms","conn":"127.0.0.1:6379","dbNum":"1","password":""}`
)

func SmsCodeServiceOf() SmsCodeService {
		var service = new(SmsCodeServiceAliCloudImpl)
		service.Init()
		return service
}

func (this *SmsCodeServiceAliCloudImpl) Init() {
		this.init()
		this.client = this.CreateClient()
		this.storage = this.initStorage()
		this.mock = getSmsCodeMockService()
		this.mock.SetContext(this)
}

func (this *SmsCodeServiceAliCloudImpl) initStorage() cache.Cache {
		if this.storage != nil {
				return this.storage
		}
		var err error
		driver := beego.AppConfig.DefaultString(SmsCacheDriverKey, SmsCacheDriverDefault)
		config := beego.AppConfig.DefaultString(SmsCacheConfigKey, SmsCacheConfigDefault)
		this.storage, err = cache.NewCache(driver, config)
		if err != nil {
				logs.Error(err)
		}
		return this.storage
}

func (this *SmsCodeServiceAliCloudImpl) CreateClient() *dysmsapi.Client {
		regionId := beego.AppConfig.DefaultString(DySmsApiRegionId, DefaultDySmsApiRegionId)
		accessKeyId, accessKeySecret := beego.AppConfig.String(DySmsAccessKeyId), beego.AppConfig.String(DySmsAccessKeySecret)
		client, err := dysmsapi.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)
		if err == nil {
				return client
		}
		return nil
}

func (this *SmsCodeServiceAliCloudImpl) Send(mobile string, content string, extras map[string]string) error {
		if this.Debug() {
				return this.mock.Send(mobile, content, extras)
		}
		extras["content"] = content
		req := this.CreateSmsRequest(mobile, extras)
		rep, err := this.client.SendSms(req)
		this.dispatch(mobile, map[string]interface{}{
				"response": rep,
				"mobile":   mobile,
				"extras":   extras,
		})

		log := new(models.SmsLog)
		log.Defaults()
		log.Content = content
		if rep != nil {
				log.Result = rep.String()
		}
		log.Mobile = mobile
		this.addLog(log)
		return err
}

func (this *SmsCodeServiceAliCloudImpl) SendCode(mobile string, typ string, extras map[string]string) (string, error) {
		if this.Debug() {
				return this.mock.SendCode(mobile, typ, extras)
		}
		var (
				code = libs.RandomNumLimitN(this.getCodeLimit())
		)
		extras["type"] = typ
		return code, this.Send(mobile, code, extras)
}

func (this *SmsCodeServiceAliCloudImpl) getCodeLimit() int {
		if count := os.Getenv("sms_code_count"); count != "" {
				if n, err := strconv.Atoi(count); err == nil && n > 0 {
						return n
				}
		}
		return beego.AppConfig.DefaultInt("sms_code_count", 6)
}

func (this *SmsCodeServiceAliCloudImpl) dispatch(mobile string, data map[string]interface{}) {

}

func (this *SmsCodeServiceAliCloudImpl) CreateSmsRequest(mobile string, extras map[string]string) *dysmsapi.SendSmsRequest {
		req := dysmsapi.CreateSendSmsRequest()
		req.Scheme = "https"
		req.PhoneNumbers = mobile
		req.SignName = extras["sign_name"]
		req.TemplateCode = extras["template_code"]
		delete(extras, "sign_name")
		delete(extras, "template_code")
		req.FormParams = extras
		return req
}

func (this *SmsCodeServiceAliCloudImpl) Verify(mobile string, code string, typ string) bool {
		data := this.Get(mobile, typ)
		if code == "" || data == "" {
				return false
		}
		if data == strings.TrimSpace(code) {
				_ = this.GetCache().Delete(this.Key(mobile, typ))
				return true
		}
		return false
}

func (this *SmsCodeServiceAliCloudImpl) GetCache() cache.Cache {
		if this.storage != nil {
				return this.storage
		}
		this.storage = this.initStorage()
		return this.storage
}

func (this *SmsCodeServiceAliCloudImpl) Storage(mobile string, data string, typ string, timeout time.Duration) {
		var key = this.Key(mobile, typ)
		_ = this.GetCache().Put(key, data, timeout)
}

func (this *SmsCodeServiceAliCloudImpl) Get(mobile string, ty string) string {
		var key = this.Key(mobile, ty)
		v := this.GetCache().Get(key)
		if v == nil {
				return ""
		}
		if str, ok := v.(string); ok {
				return str
		}
		if d, ok := v.([]byte); ok {
				return string(d)
		}
		return ""
}

func (this *SmsCodeServiceAliCloudImpl) Key(mobile string, ty string) string {
		return fmt.Sprintf("%s:%s", mobile, ty)
}

func (this *SmsCodeServiceAliCloudImpl) addLog(log *models.SmsLog) {
		go func(log *models.SmsLog) {
				_ = models.SmsLogModelOf().Add(log)
		}(log)
}

func (this *SmsCodeServiceAliCloudImpl) Debug() bool {
		return beego.BConfig.RunMode == "dev" || os.Getenv("sms_debug_on") == "1"
}

type smsCodeMockService struct {
		config map[string]string
		locker sync.Mutex
		caller SmsCodeService
}

var (
		mockLocker     sync.Once
		mockSmsService *smsCodeMockService
)

func getSmsCodeMockService() *smsCodeMockService {
		if mockSmsService == nil {
				mockLocker.Do(newMockSmsService)
		}
		return mockSmsService
}

func newMockSmsService() {
		mockSmsService = new(smsCodeMockService)
		mockSmsService.init()
}

func (this *smsCodeMockService) init() {
		this.locker = sync.Mutex{}
		this.config = make(map[string]string)
}

func (this *smsCodeMockService) Send(mobile string, content string, extras map[string]string) error {
		logs.Debug(fmt.Sprintf("mobile:%s, content:%s extras:%v \n", mobile, content, extras))
		return nil
}

func (this *smsCodeMockService) SendCode(mobile string, typ string, extras map[string]string) (string, error) {
		var code = this.getCode()
		this.Storage(mobile, code, typ, 3*time.Minute)
		return code, nil
}

func (this *smsCodeMockService) getCode() string {
		if code := os.Getenv("sms_mock_code"); code != "" {
				return code
		}
		return "111111"
}

func (this *smsCodeMockService) Verify(mobile string, code string, typ string) bool {
		return this.caller.Verify(mobile, code, typ)
}

func (this *smsCodeMockService) Storage(mobile string, data string, typ string, timeout time.Duration) {
		this.caller.Storage(mobile, data, typ, timeout)
}

func (this *smsCodeMockService) SetContext(ctx SmsCodeService) {
		this.caller = ctx
}
