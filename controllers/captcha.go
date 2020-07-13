package controllers

import "github.com/weblfe/travel-app/repositories"

type CaptchaController struct {
	BaseController
}

// 验证码模块 controller
func CaptchaControllerOf() *CaptchaController  {
	 return new(CaptchaController)
}

// @router /captcha/mobile [post]
func (this *CaptchaController)SendMobileCaptcha()  {
   this.Send(repositories.NewCaptchaRepository(this).SendMobileCaptcha())
}

// @router /captcha/email [post]
func (this *CaptchaController)SendEmailCaptcha()  {

}

// @router /captcha/wechat [post]
func (this *CaptchaController)SendWeChatCaptcha()  {

}