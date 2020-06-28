package repositories

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/services"
		"io"
		"os"
		"path"
		"time"
)

type AttachmentRepository interface {
		Ticket() common.ResponseJson
		Upload() common.ResponseJson
		Uploads() common.ResponseJson
		List() common.ResponseJson
}

type AttachmentRepositoryImpl struct {
		ctx        *beego.Controller
		filesystem services.FileSystem
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
		this.filesystem = services.GetFileSystem()
}

func (this *AttachmentRepositoryImpl) List() common.ResponseJson {

		return common.NewInvalidParametersResp()
}

func (this *AttachmentRepositoryImpl) Ticket() common.ResponseJson {
		panic("implement me")
}

func (this *AttachmentRepositoryImpl) Upload() common.ResponseJson {
		var (
				ctx = this.ctx
				typ = ctx.GetString("type")
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
		dir := path.Join(beego.AppPath, "static/storage")
		if !libs.IsExits(dir) {
				_ = os.MkdirAll(dir, os.ModePerm)
		}
		filename := path.Join(dir, fs.Filename)
		fd, errs := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
		if errs != nil {
				return common.NewErrorResp(common.NewErrors(common.ServiceFailed, errs.Error()))
		}
		defer func() {
				_ = fd.Close()
		}()
		_, err = io.Copy(fd, m)
		if err != nil {
				return common.NewErrorResp(common.NewErrors(common.ServiceFailed, "文件保存失败"))
		}
		id := libs.FileHash(filename)
		return common.NewSuccessResp(beego.M{"mediaId": id, "time": time.Now().Unix()}, "上传成功")
}

func (this *AttachmentRepositoryImpl) Uploads() common.ResponseJson {
		panic("implement me")
}
