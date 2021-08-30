package controllers

import "github.com/weblfe/travel-app/repositories"

type CaptchaController struct {
	BaseController
}

// CaptchaControllerOf 验证码模块 controller
func CaptchaControllerOf() *CaptchaController  {
	 return new(CaptchaController)
}

// SendMobileCaptcha
// @router /captcha/mobile [post]
func (this *CaptchaController)SendMobileCaptcha()  {
   this.Send(repositories.NewCaptchaRepository(this).SendMobileCaptcha())
}

// SendEmailCaptcha
// @router /captcha/email [post]
func (this *CaptchaController)SendEmailCaptcha()  {

}

// SendWeChatCaptcha
// @router /captcha/wechat [post]
func (this *CaptchaController)SendWeChatCaptcha()  {

}