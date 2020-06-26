package services

import (
		"fmt"
		"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/cache"
		"github.com/weblfe/travel-app/models"
		"math/rand"
		"time"
)

type SmsCodeService interface {
		Send(mobile string, content string, extras map[string]string) error
		SendCode(mobile string, typ string, extras map[string]string) (string, error)
		Verify(mobile string, code string, typ string) bool
		Storage(mobile string, data string, typ string, timeout time.Duration)
}

type SmsCodeServiceAliCloudImpl struct {
		client  *dysmsapi.Client
		storage cache.Cache
}

const (
		DySmsApiRegionId        = "dysms_api_region_id"
		DySmsAccessKeyId        = "dysms_access_key_id"
		DySmsAccessKeySecret    = "dysms_access_key_secret"
		DefaultDySmsApiRegionId = "cn-hangzhou"
		SmsCacheDriverKey       = "sms_cache_driver"
		SmsCacheDriverDefault   = "redis"
		SmsCacheConfigKey       = "sms_cache_config"
		SmsCacheConfigDefault   = `{"key":"sms","conn":"127.0.0.1:6039","dbNum":"1","password":""}`
)

func SmsCodeServiceOf() SmsCodeService {
		var service = new(SmsCodeServiceAliCloudImpl)
		service.init()
		return service
}

func (this *SmsCodeServiceAliCloudImpl) init() {
		this.client = this.CreateClient()
		this.storage = this.initStorage()
}

func (this *SmsCodeServiceAliCloudImpl) initStorage() cache.Cache {
		if this.storage != nil {
				return this.storage
		}
		driver := beego.AppConfig.DefaultString(SmsCacheDriverKey, SmsCacheDriverDefault)
		config := beego.AppConfig.DefaultString(SmsCacheConfigKey, SmsCacheConfigDefault)
		this.storage, _ = cache.NewCache(driver, config)
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
		var (
				code = fmt.Sprintf("%d", rand.Intn(6))
		)
		extras["type"] = typ
		return code, this.Send(mobile, code, extras)
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

		return false
}

func (this *SmsCodeServiceAliCloudImpl) GetCache() cache.Cache {
		if this.storage == nil {
				return this.storage
		}
		this.storage = this.initStorage()
		return this.storage
}

func (this *SmsCodeServiceAliCloudImpl) Storage(mobile string, data string, typ string, timeout time.Duration) {
		var key = fmt.Sprintf("%s:%s", mobile, typ)
		_ = this.GetCache().Put(key, data, timeout)
}

func (this *SmsCodeServiceAliCloudImpl) Get(mobile string, ty string) string {
		var key = fmt.Sprintf("%s:%s", mobile, ty)
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

func (this *SmsCodeServiceAliCloudImpl) addLog(log *models.SmsLog) {
		go func(log *models.SmsLog) {
				_ = models.SmsLogModelOf().Add(log)
		}(log)
}
