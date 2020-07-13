package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
		"github.com/weblfe/travel-app/transports"
		"time"
)

type CaptchaRepository interface {
		SendMobileCaptcha() common.ResponseJson
}

type CaptchaRepositoryImpl struct {
		smsService services.SmsCodeService
		ctx        common.BaseRequestContext
}

func NewCaptchaRepository(ctx common.BaseRequestContext) CaptchaRepository {
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
				request = new(transports.MobileRequest)
		)
		request.Load(this.ctx.GetParent().Ctx.Input)
		request.Dump()
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
