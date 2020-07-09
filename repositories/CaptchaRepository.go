package repositories

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
		"time"
)

type CaptchaRepository interface {
		SendMobileCaptcha() common.ResponseJson
}

type CaptchaRepositoryImpl struct {
		smsService services.SmsCodeService
		ctx        *beego.Controller
}

func NewCaptchaRepository(ctx *beego.Controller) CaptchaRepository {
		var repository = new(CaptchaRepositoryImpl)
		repository.ctx = ctx
		repository.init()
		return repository
}

func (this *CaptchaRepositoryImpl) init() {
		this.smsService = services.SmsCodeServiceOf()
}

func (this *CaptchaRepositoryImpl) SendMobileCaptcha() common.ResponseJson {
		var (
				mobile  string
				typ     string
				ctx     = this.ctx.Ctx
				request = new(MobileRequest)
		)
		request.Load(ctx.Input)
		typ = request.Type
		mobile = request.Mobile
		if mobile == "" || typ == "" {
				return common.NewInvalidParametersResp(common.NewErrors(common.InvalidParametersCode, "手机号和发送类型不能为空"))
		}
		code, err := this.smsService.SendCode(mobile, typ, map[string]string{
				"type": "sms_verify_code",
		})
		if err != nil {
				return common.NewErrorResp(common.NewErrors(common.ServiceFailed, err.Error()))
		}
		this.smsService.Storage(mobile, code, typ, 6*time.Minute)
		return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix()}, "发送成功,请在6分钟内使用")
}

// 验证码参数
type MobileRequest struct {
		Mobile string `json:"mobile"`
		Type   string `json:"type"`
}

// 加载
func (this *MobileRequest) Load(ctx *context.BeegoInput) *MobileRequest {
		if err := ctx.Bind(&this.Mobile, "mobile"); err != nil || this.Mobile == "" {
				_ = json.Unmarshal(ctx.RequestBody, this)
		}
		if this.Type == "" {
				_ = ctx.Bind(&this.Type, "type")
		}
		return this
}

// 过滤
func (this *MobileRequest) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = beego.M{
				"mobile": this.Mobile,
				"type" : this.Type,
		}
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}
