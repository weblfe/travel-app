package services

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"io"
)

type AttachmentService interface {
		Remove(query beego.M) bool
		Get(mediaId string) *models.Attachment
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
		if model == nil || !this.Create(model) {
				return nil
		}
		return model
}

func (this *AttachmentServiceImpl) Create(attach *models.Attachment) bool {
		if err := this.model.Add(attach); err == nil {
				return false
		}
		return true
}

func (this *AttachmentServiceImpl) defaultsExtras(m beego.M) beego.M {
		_, ok := m["path"]
		if len(m) == 0 || !ok {
				m["path"] = PathsServiceOf().StoragePath()
		}
		return m
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
				res, ok := GetFileSystem().SaveByReader(reader, extras)
				if ok && len(res) > 0 {
						return models.NewAttachment().Load(res).Defaults()
				}
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
