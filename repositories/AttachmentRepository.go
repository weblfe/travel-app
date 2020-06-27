package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
)

type AttachmentRepository interface {
		Ticket() common.ResponseJson
		Upload() common.ResponseJson
		Uploads() common.ResponseJson
		List() common.ResponseJson
}

type AttachmentRepositoryImpl struct {
		ctx *beego.Controller
}

func NewAttachmentRepository(ctx *beego.Controller) AttachmentRepository  {
		var repository = new(AttachmentRepositoryImpl)
		repository.ctx = ctx
		return repository
}

func (this *AttachmentRepositoryImpl)List() common.ResponseJson  {

		return common.NewInvalidParametersResp()
}

func (this *AttachmentRepositoryImpl) Ticket() common.ResponseJson {
		panic("implement me")
}

func (this *AttachmentRepositoryImpl) Upload() common.ResponseJson {
		panic("implement me")
}

func (this *AttachmentRepositoryImpl) Uploads() common.ResponseJson {
		panic("implement me")
}