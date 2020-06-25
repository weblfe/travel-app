package controllers

import "github.com/astaxie/beego"

type CaptchaController struct {
		beego.Controller
}

// 验证码模块 controller
func CaptchaControllerOf() *CaptchaController  {
	 return new(CaptchaController)
}

// @route /captcha/mobile [post]
func (this *CaptchaController)SendMobileCaptcha()  {

}

// @route /captcha/email [post]
func (this *CaptchaController)SendEmailCaptcha()  {

}

// @route /captcha/wechat [post]
func (this *CaptchaController)SendWeChatCaptcha()  {

}