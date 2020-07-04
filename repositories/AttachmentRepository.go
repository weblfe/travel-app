package repositories

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
		"io"
		"net/http"
		"path/filepath"
)

type AttachmentRepository interface {
		Ticket() common.ResponseJson
		Upload() common.ResponseJson
		Uploads() common.ResponseJson
		List() common.ResponseJson
		GetByMediaId(...string)
		DownloadByMediaId(...string)
}

type AttachmentRepositoryImpl struct {
		ctx               *beego.Controller
		attachmentService services.AttachmentService
}

const (
		DefaultFileKey = "file"
)

func NewAttachmentRepository(ctx *beego.Controller) AttachmentRepository {
		var repository = new(AttachmentRepositoryImpl)
		repository.ctx = ctx
		repository.init()
		return repository
}

func (this *AttachmentRepositoryImpl) init() {
		this.attachmentService = services.AttachmentServiceOf()
}

func (this *AttachmentRepositoryImpl) List() common.ResponseJson {

		return common.NewInvalidParametersResp()
}

func (this *AttachmentRepositoryImpl) Ticket() common.ResponseJson {
		panic("implement me")
}

func (this *AttachmentRepositoryImpl) Upload() common.ResponseJson {
		var (
				ctx    = this.ctx
				typ    = ctx.GetString("type")
				ticket = ctx.GetString("ticket")
		)
		if typ == "" {
				typ = DefaultFileKey
		}
		m, fs, err := ctx.GetFile(typ)
		if m != nil {
				defer func(closer io.Closer) {
						err := closer.Close()
						if err != nil {
								logs.Error(err)
						}
				}(m)
		}
		if err != nil {
				return common.NewErrorResp(common.NewErrors(common.InvalidTokenCode, err.Error()))
		}
		if m == nil {
				return common.NewErrorResp(common.NewErrors(common.InvalidTokenCode, "文件异常"))
		}
		if fs == nil {
				return common.NewErrorResp(common.NewErrors(common.InvalidTokenCode, "文件传输异常"))
		}
		// 保存
		result := this.attachmentService.Save(m, beego.M{"fileInfo": fs, "ticket": ticket, "path": services.PathsServiceOf().StoragePath()})
		if result == nil {
				return common.NewErrorResp(common.NewErrors(common.InvalidTokenCode, "文件保存失败"))
		}
		return common.NewSuccessResp(result.M(FilterWrapper(filterAttachment, filterEmptyMapper)), "上传成功")
}

func (this *AttachmentRepositoryImpl) Uploads() common.ResponseJson {
		panic("implement me")
}

// 下载文件
func (this *AttachmentRepositoryImpl) DownloadByMediaId(mediaIds...string) {
		var (
				id   = this.ctx.GetString(":mediaId")
				info = this.attachmentService.Get(id)
		)
		if id == "" && len(mediaIds) > 0 {
				id = mediaIds[0]
		}
		if info == nil || id == ""{
				this.ctx.Ctx.Output.Status = 404
				this.ctx.Ctx.WriteString("media file not found!")
				return
		}
		if info.Path == "" {
				this.ctx.Ctx.Output.Status = 404
				this.ctx.Ctx.WriteString("media file not found!")
				return
		}
		this.ctx.Ctx.Output.Download(filepath.Join(info.Path, info.FileName), id+"."+filepath.Ext(info.FileName))
		return
}

// 文件服务
func (this *AttachmentRepositoryImpl) GetByMediaId(mediaIds...string) {
		var (
				id   = this.ctx.GetString(":mediaId")
				info = this.attachmentService.Get(id)
		)
		if id == "" && len(mediaIds) > 0 {
				id = mediaIds[0]
		}
		if info == nil || id == ""{
				this.ctx.Ctx.Output.Status = 404
				this.ctx.Ctx.WriteString("media file not found!")
				return
		}
		if info.Path == "" {
				this.ctx.Ctx.Output.Status = 404
				this.ctx.Ctx.WriteString("media file not found!")
				return
		}
		http.ServeFile(this.ctx.Ctx.ResponseWriter,this.ctx.Ctx.Request,filepath.Join(info.Path, info.FileName))
		return
}