package controllers

type PopularizationController struct {
		BaseController
}

// 推广 控制器
func PopularizationControllerOf() *PopularizationController {
		return new(PopularizationController)
}

// 服务 渠道信息
func (this *PopularizationController)GetChannelInfo()  {

}

// 获取 推广
func (this *PopularizationController)GetChannelQrCode()  {

}

// 获取 推广码
func (this *PopularizationController)GetChannelCode()  {

}