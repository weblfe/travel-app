package controllers

import (
		"github.com/weblfe/travel-app/repositories"
)

type AttachmentController struct {
		BaseController
}

func AttachmentControllerOf() *AttachmentController {
		return new(AttachmentController)
}

// 上传令牌
// @router /attachment/ticket [post]
func (this *AttachmentController) Ticket() {
		this.Send(repositories.NewAttachmentRepository(this).Ticket())
}

// 上传单个文件
// @router /attachment/upload [post]
func (this *AttachmentController) Upload() {
		this.Send(repositories.NewAttachmentRepository(this).Upload())
}

// 上传多文件
// @router /attachment/uploads [post]
func (this *AttachmentController) Uploads() {
		this.Send(repositories.NewAttachmentRepository(this).Uploads())
}

// 获取附件列表
// @router /attachment/list [get]
func (this *AttachmentController) List() {
		this.Send(repositories.NewAttachmentRepository(this).List())
}

// 下载附件
// @router /attachments/download/:mediaId  [get]
func (this *AttachmentController) DownloadByMediaId() {
		repositories.NewAttachmentRepository(this).DownloadByMediaId()
}

// 浏览附件
// @router /attachments/:mediaId  [get]
func (this *AttachmentController) GetByMediaId() {
		repositories.NewAttachmentRepository(this).GetByMediaId()
}
