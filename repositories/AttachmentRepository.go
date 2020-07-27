package repositories

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/services"
		"github.com/weblfe/travel-app/transforms"
		"io"
		"net/http"
		"path/filepath"
		"time"
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
		ctx               common.BaseRequestContext
		attachmentService services.AttachmentService
		urlTicketService  services.UrlTicketService
}

const (
		DefaultFileKey  = "file"
		DefaultFilesKey = "files"
)

func NewAttachmentRepository(ctx common.BaseRequestContext) AttachmentRepository {
		var repository = new(AttachmentRepositoryImpl)
		repository.ctx = ctx
		repository.init()
		return repository
}

func (this *AttachmentRepositoryImpl) init() {
		this.attachmentService = services.AttachmentServiceOf()
}

// 罗列附件接口
func (this *AttachmentRepositoryImpl) List() common.ResponseJson {
		return common.NewInDevResp(this.ctx.GetActionId())
}

// 跨站上传ticket
func (this *AttachmentRepositoryImpl) Ticket() common.ResponseJson {
		return common.NewInDevResp(this.ctx.GetActionId())
}

// 单文件上传
func (this *AttachmentRepositoryImpl) Upload() common.ResponseJson {
		var (
				ctx      = this.ctx.GetParent()
				typ      = ctx.GetString("type")
				ticket   = ctx.GetString("ticket")
				uid      = ctx.Ctx.Input.Param("_userId")
				fileType = ctx.GetString("fileType")
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
		data := beego.M{
				"userId": uid, "fileInfo": fs,
				"ticket": ticket, "path": this.getAttachmentPath(),
				"type": fileType,
		}
		result := this.attachmentService.Save(m, data)
		if result == nil {
				return common.NewErrorResp(common.NewErrors(common.InvalidTokenCode, "文件保存失败"))
		}
		filter := transforms.FilterWrapper(this.getAttachmentFilters()...)
		data  = result.M(filter)
		data["url"] = services.UrlTicketServiceOf().GetTicketUrlByAttach(result)
		return common.NewSuccessResp(data, "上传成功")
}

// 多文件上传
func (this *AttachmentRepositoryImpl) Uploads() common.ResponseJson {
		var (
				ctx      = this.ctx.GetParent()
				typ      = ctx.GetString("type")
				ticket   = ctx.GetString("ticket")
				uid      = ctx.Ctx.Input.Param("_userId")
				fileType = ctx.Ctx.Input.Param("fileType")
		)

		if typ == "" {
				typ = DefaultFilesKey
		}
		fsArr, err := ctx.GetFiles(typ)
		if err != nil {
				return common.NewErrorResp(common.NewErrors(common.ServiceFailed, err.Error()))
		}
		var (
				failCount    int
				successCount int
				results      []beego.M
				filter       = transforms.FilterWrapper(this.getAttachmentFilters()...)
		)
		// 是否同一中类型
		if fileType != "" {
				tmp := ""
				for _, fs := range fsArr {
						cur := libs.GetFileType(fs.Filename)
						if tmp != "" {
								tmp = cur
								continue
						}
						if tmp != cur {
								fileType = ""
						}
				}
		}
		var ticketService = services.UrlTicketServiceOf()
		// 文件保存
		for _, m := range fsArr {
				fs, openErr := m.Open()
				if openErr != nil {
						results = append(results, beego.M{"filename": m.Filename, "status": 0, "error": openErr.Error()})
						continue
				}
				data := beego.M{
						"userId": uid, "fileInfo": m,
						"ticket": ticket, "path": this.getAttachmentPath(),
				}
				if fileType != "" {
						data["type"] = fileType
				}
				// 保存
				result := this.attachmentService.Save(fs, data)
				_ = fs.Close()
				if result == nil {
						failCount++
						results = append(results, beego.M{"filename": m.Filename, "status": -1, "error": "save failed!"})
				} else {
						successCount++
						it:= result.M(filter)
						it["url"] = ticketService.GetTicketUrlByAttach(result)
						results = append(results, it)
				}
		}
		if successCount == 0 {
				return common.NewErrorResp(common.NewErrors("all save failed!", common.ServiceFailed), "文保存失败")
		}
		data := beego.M{
				"items": results,
				"meta":  beego.M{"successCount": successCount, "failCount": failCount, "total": len(fsArr)},
		}
		return common.NewSuccessResp(data, "上传成功")
}

// 下载文件
func (this *AttachmentRepositoryImpl) DownloadByMediaId(mediaIds ...string) {
		var (
				ctx  = this.ctx.GetParent()
				id   = ctx.GetString(":mediaId")
				info = this.attachmentService.Get(id)
		)
		if id == "" && len(mediaIds) > 0 {
				id = mediaIds[0]
		}
		if info == nil || id == "" {
				ctx.Ctx.Output.Status = 404
				ctx.Ctx.WriteString("media file not found!")
				return
		}
		if info.Path == "" {
				ctx.Ctx.Output.Status = 404
				ctx.Ctx.WriteString("media file not found!")
				return
		}
		ctx.Ctx.Output.Download(filepath.Join(info.Path, info.FileName), id+"."+filepath.Ext(info.FileName))
		return
}

// 文件服务
func (this *AttachmentRepositoryImpl) GetByMediaId(mediaIds ...string) {
		var (
				ctx  = this.ctx.GetParent()
				id   = ctx.GetString(":mediaId")
				info = this.attachmentService.Get(id)
		)
		if id == "" && len(mediaIds) > 0 {
				id = mediaIds[0]
		}
		if info == nil || id == "" {
				ctx.Ctx.Output.Status = 404
				ctx.Ctx.WriteString("media file not found!")
				return
		}
		if info.Path == "" {
				ctx.Ctx.Output.Status = 404
				ctx.Ctx.WriteString("media file not found!")
				return
		}
		http.ServeFile(ctx.Ctx.ResponseWriter, ctx.Ctx.Request, filepath.Join(info.Path, info.FileName))
		return
}

func (this *AttachmentRepositoryImpl) getAttachmentPath() string {
		var (
				year, month, day = time.Now().Date()
				date             = fmt.Sprintf("%d-%d-%d", year, month, day)
		)
		return services.PathsServiceOf().StoragePath("/" + date)
}

func (this *AttachmentRepositoryImpl) getAttachmentFilters() []func(beego.M) beego.M {
		return []func(beego.M) beego.M{
				transforms.FilterAttachment, transforms.FilterEmptyMapper,
				transforms.FieldsFilter([]string{"path", "id", "createdAt","updatedAt", "extrasInfo"}),
		}
}

func (this *AttachmentRepositoryImpl) transUrl(m beego.M) beego.M {
		var id = m["mediaId"]
		if id == nil || id == "" {
				return m
		}
		m["url"] = this.attachmentService.GetUrl(id.(string))
		return m
}
