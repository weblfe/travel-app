package services

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config/env"
	"github.com/astaxie/beego/logs"
	"github.com/weblfe/travel-app/common"
	"github.com/weblfe/travel-app/libs"
	"github.com/weblfe/travel-app/models"
	"io"
	"os"
	"path/filepath"
	"time"
)

type AttachmentService interface {
	Remove(query beego.M) bool
	Get(mediaId string) *models.Attachment
	UpdateById(string, beego.M) error
	GetUrl(string) string
	Save(reader io.ReadCloser, extras ...beego.M) *models.Attachment
}

type AttachmentServiceImpl struct {
	BaseService
	model *models.AttachmentModel
}

func AttachmentServiceOf() AttachmentService {
	var service = new(AttachmentServiceImpl)
	service.Init()
	return service
}

func (this *AttachmentServiceImpl) Init() {
	this.model = models.AttachmentModelOf()
	this.Constructor = func(args ...interface{}) interface{} {
		return AttachmentServiceOf()
	}
	this.init()
}

func (this *AttachmentServiceImpl) Remove(query beego.M) bool {
	if err := this.model.Remove(query, true); err == nil {
		return true
	}
	return false
}

func (this *AttachmentServiceImpl) Get(mediaId string) *models.Attachment {
	var model, err = this.model.GetByMediaId(mediaId)
	if err == nil {
		return model
	}
	return nil
}

func (this *AttachmentServiceImpl) Save(reader io.ReadCloser, extras ...beego.M) *models.Attachment {
	var m = this.defaultsExtras(libs.MapMerge(extras...))
	if reader == nil {
		return nil
	}
	model := this.save(reader, m)
	if model == nil {
		return nil
	}
	if ty, ok := m["type"]; ok && ty != nil && ty != "" {
		tyName := ty.(string)
		if tyName != "file" && tyName != "files" {
			model.FileType = tyName
		}
	}
	if this.Create(model) {
		return model
	}
	return nil
}

func (this *AttachmentServiceImpl) Create(attach *models.Attachment) bool {
	if attach == nil {
		return false
	}
	attach = this.onlySaveOne(attach)
	if err := this.model.Add(attach); err == nil {
		return true
	}
	return false
}

func (this *AttachmentServiceImpl) delete(fs string) {
	if fs != "" {
		err := os.Remove(fs)
		if err != nil {
			logs.Error(err)
		}
	}
}

// 文件仅保存一份
func (this *AttachmentServiceImpl) onlySaveOne(attach *models.Attachment) *models.Attachment {
	// 开关文件保存一份 通过hash
	if env.Get("ATTACHMENT_ONLY_ONE_SAVE", "on") == "off" {
		return attach
	}
	if attach.Hash != "" {
		oldAttach := this.GetByHash(attach.Hash)
		if oldAttach != nil && attach.Path != "" && attach.FileName != "" {
			existsFs := filepath.Join(oldAttach.Path, oldAttach.FileName)
			if !libs.IsExits(existsFs) {
				return attach
			}
			fs := filepath.Join(attach.Path, attach.FileName)
			defer this.delete(fs)
			attach.ExtrasInfo["originSavePath"] = attach.Path
			attach.ExtrasInfo["originSaveFileName"] = attach.FileName
			attach.Path = oldAttach.Path
			attach.FileName = oldAttach.FileName
		}
	}
	return attach
}

func (this *AttachmentServiceImpl) defaultsExtras(m beego.M) beego.M {
	_, ok := m["path"]
	if len(m) == 0 || !ok {
		m["path"] = this.getAttachmentPath()
	}
	return m
}

func (this *AttachmentServiceImpl) getAttachmentPath() string {
	var (
		year, month, day = time.Now().Date()
		date             = fmt.Sprintf("%d-%d-%d", year, month, day)
	)
	return PathsServiceOf().StoragePath("/" + date)
}

func (this *AttachmentServiceImpl) save(reader io.ReadCloser, extras beego.M) *models.Attachment {
	var (
		path           = extras["path"]
		oss, ossBucket = extras["oss"], extras["oss_bucket"]
	)
	if reader == nil {
		return nil
	}
	if path != "" {
		extras["path"] = path
		_ = os.MkdirAll(path.(string), os.ModePerm)
		res, ok := GetFileSystem().SaveByReader(reader, extras)
		if ok && len(res) > 0 {
			return models.NewAttachment().Load(libs.MapMerge(res, extras)).Defaults()
		}
		return nil
	}
	if oss != "" && ossBucket != "" {
		extras["schema"] = "@oss(" + oss.(string) + ")://" + ossBucket.(string)
		res, ok := GetFileSystem().SaveByReader(reader, extras)
		if ok && len(res) > 0 {
			return models.NewAttachment().Load(res).Defaults()
		}
	}
	return nil
}

func (this *AttachmentServiceImpl) GetByHash(hash string) *models.Attachment {
	var attach = models.NewAttachment()
	if err := this.model.GetByKey("hash", hash, attach); err == nil {
		return attach
	}
	return nil
}

func (this *AttachmentServiceImpl) UpdateById(id string, update beego.M) error {
	if len(update) == 0 || id == "" {
		return common.NewErrors(common.InvalidParametersCode, "更新参数不能为空")
	}
	return this.model.UpdateById(id, update)
}

func (this *AttachmentServiceImpl) GetUrl(mediaId string) string {
	var data = this.Get(mediaId)
	if data == nil {
		return ""
	}
	if data.CdnUrl != "" {
		return data.CdnUrl
	}
	return data.Url
}
