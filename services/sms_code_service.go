package services

import (
		"encoding/json"
		"fmt"
		"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/cache"
		_ "github.com/astaxie/beego/cache/memcache"
		_ "github.com/astaxie/beego/cache/redis"
		"github.com/astaxie/beego/config/env"
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo/bson"
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

// Send 发送短信
func (this *SmsCodeServiceAliCloudImpl) Send(mobile string, content string, extras map[string]string) error {
		if this.Debug() {
				return this.mock.Send(mobile, content, extras)
		}
		// 填写消息内容
		extras["content"] = content
		typ  :=extras["type"]
		req := this.CreateSmsRequest(mobile, extras)
		if req == nil {
				return fmt.Errorf("参数不足 %v", extras)
		}
		// sdk 发送
		rep, err := this.client.SendSms(req)
		// 派发结果
		defer this.dispatch(mobile, map[string]interface{}{
				"response": rep,
				"mobile":   mobile,
				"extras":   extras,
				"request" : req,
		})

		log := new(models.SmsLog)
		log.Defaults()
		log.Content = content
		log.Type = typ
		tmp, _ := json.Marshal(extras)
		log.Extras = string(tmp)
		if rep != nil {
				log.Result = rep.GetHttpContentString()
		}
		if err !=nil {
       log.State = 2
       log.Result = "error: "+ err.Error() + " response: " + log.Result
		}else{
				log.State = 1
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
		extras["name"] = typ
		extras["code"] = code
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
		logs.Debug("mobile :" + mobile + fmt.Sprintf("result: %v", data))
}

// CreateSmsRequest 创建请求体
func (this *SmsCodeServiceAliCloudImpl) CreateSmsRequest(mobile string, extras map[string]string) *dysmsapi.SendSmsRequest {
		req := dysmsapi.CreateSendSmsRequest()
		req.PhoneNumbers = mobile
		req.Scheme = this.getScheme(extras)
		req.SignName = this.getSignName(extras)
		req.TemplateCode = this.getTemplateCode(&extras)
		this.cleanExtras(extras)
		if len(extras) > 0 {
				req.FormParams = extras
				tmp,_ := json.Marshal(extras)
				if len(tmp) >0 {
						req.TemplateParam = string(tmp)
				}
		}
		if req.TemplateCode == "" || len(req.FormParams) <= 0 {
				return nil
		}
		return req
}

func (this *SmsCodeServiceAliCloudImpl) cleanExtras(extras map[string]string) map[string]string {
		delete(extras, "sign_name")
		delete(extras, "template_code")
		delete(extras, "type")
		delete(extras, "name")
		delete(extras, "platform")
		delete(extras, "scheme")
		delete(extras, "content")
		return extras
}

func (this *SmsCodeServiceAliCloudImpl) getScheme(extras map[string]string) string {
		if v, ok := extras["scheme"]; ok {
				if libs.InArray(v, []string{"http", "https"}) {
						return v
				}
		}
		return env.Get("DYSMS_SCHEME", "https")
}

func (this *SmsCodeServiceAliCloudImpl) getSignName(extras map[string]string) string {
		if v, ok := extras["sign_name"]; ok && v != "" {
				return v
		}
		return env.Get("DYSMS_SIGN_NAME", "绿游App")
}

// 设置 模版code 和相关请求参数
func (this *SmsCodeServiceAliCloudImpl) getTemplateCode(extras *map[string]string) string {
		var (
				_extras     = *extras
				name, _     = _extras["name"]
				typ, _      = _extras["type"]
				platform, _ = _extras["platform"]
		)
		if platform == "" {
				platform = "ali_dy"
		}
		if name == ""  {
				if v, ok := _extras["template_code"]; ok && v != "" {
						return v
				}
				return ""
		}
		// @todo 常量 默认验证码类型
		if  typ == "" {
				typ = "sms_verify_code"
		}
		if !strings.Contains(name,"_sms_code") {
				name += "_sms_code"
		}
		var (
				tpl   = models.NewMessageTemplate()
				model = models.MessageTemplateModelOf()
				query = bson.M{
						"name":     name,
						"type":     typ,
						"platform": platform,
				}
		)
		if err := model.FindOne(query, tpl); err == nil && tpl.TemplateId != "" {
				// 设置请求参数
				this.setFromDataByTemplate(tpl, extras)
				return tpl.TemplateId
		}
		return ""
}

// 填充请求参数
func (this *SmsCodeServiceAliCloudImpl) setFromDataByTemplate(tpl *models.MessageTemplate, extras *map[string]string) {
		var (
				_extras = *extras
				data    = tpl.Template
				smsTemp = models.NewSmsTemplate().Load(data)
		)
		for _, it := range smsTemp.Variables {
				if v, ok := _extras[it.Key]; ok {
						_extras[it.Value] = v
						continue
				}
				if v, ok := _extras[it.Value]; ok {
						_extras[it.Value] = v
						continue
				}
				if it.Value != "code" {
						continue
				}
				if v, ok := _extras["content"]; ok && v != "" {
						_extras[it.Value] = v
				}
		}
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
		return beego.BConfig.RunMode == "dev" && os.Getenv("SMS_DEBUG_ON") == "1"
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
