package controllers

import "github.com/weblfe/travel-app/repositories"

type PopularizationController struct {
		BaseController
}

// 推广 控制器
func PopularizationControllerOf() *PopularizationController {
		return new(PopularizationController)
}

// 服务 渠道信息
// @router /popularization/info [get]
func (this *PopularizationController) GetChannelInfo() {
		this.Send(repositories.NewPopularizationRepository(this).GetChannelInfo())
}

// 获取 推广 二维码
// @router /popularization/qrcode [get]
func (this *PopularizationController) GetChannelQrCode() {
		this.Send(repositories.NewPopularizationRepository(this).GetChannelQrcode())
}

// 获取 推广码
// @router /popularization/channel [get]
func (this *PopularizationController) PublishChannelCode() {
   this.Send(repositories.NewPopularizationRepository(this).GetChannel())
}
