package controllers

import "github.com/astaxie/beego"

type AttachmentController struct {
		beego.Controller
}

func AttachmentControllerOf() *AttachmentController {
		return new(AttachmentController)
}

// 上传令牌
// @route /attachment/ticket
func (this *AttachmentController) ticket() {

}

// 上传单个文件
// @route /attachment/upload
func (this *AttachmentController) Upload() {

}

// 上传多文件
// @route /attachment/uploads
func (this *AttachmentController) Uploads() {

}

// 获取附件列表
// @route /attachment/list
func (this *AttachmentController) List() {

}
