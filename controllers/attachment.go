package controllers

type AttachmentController struct {
		BaseController
}

func AttachmentControllerOf() *AttachmentController {
		return new(AttachmentController)
}

// 上传令牌
// @router /attachment/ticket
func (this *AttachmentController) Ticket() {

}

// 上传单个文件
// @router /attachment/upload
func (this *AttachmentController) Upload() {

}

// 上传多文件
// @router /attachment/uploads
func (this *AttachmentController) Uploads() {

}

// 获取附件列表
// @router /attachment/list
func (this *AttachmentController) List() {

}
