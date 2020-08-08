package services

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/config/env"
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/plugins"
		"io"
		"os"
		"path/filepath"
		"strings"
		"time"
)

type AttachmentService interface {
		Remove(query beego.M) bool
		Get(mediaId string) *models.Attachment
		UpdateById(string, beego.M) error
		GetUrl(string) string
		GetById(id string) *models.Attachment
		GetAccessUrl(string) string
		AutoCoverForVideo(attachment *models.Attachment,posts...*models.TravelNotes) string
		Save(reader io.ReadCloser, extras ...beego.M) *models.Attachment
}

type AttachmentServiceImpl struct {
		BaseService
		model *models.AttachmentModel
}

const (
		AttachTypeImage       = "image"
		AttachTypeImageAvatar = "avatar"
		AttachTypeDoc         = "doc"
		AttachTypeText        = "txt"
		AttachTypeVideo       = "video"
)

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
		var attach = UrlTicketServiceOf().GetTicketInfoToSimple(mediaId)
		if attach == nil {
				var model, err = this.model.GetByMediaId(mediaId)
				if err == nil {
						return model
				}
				return nil
		}
		if model, err := this.model.GetByMediaId(attach.MediaId); model != nil && err == nil {
				return model
		}
		return nil
}

func (this *AttachmentServiceImpl) GetById(id string) *models.Attachment {
		var model, err = this.model.GetByMediaId(id)
		if model != nil && err == nil {
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
				go this.after(attach)
				return true
		}
		return false
}

func (this *AttachmentServiceImpl) after(attachment *models.Attachment) {
		switch attachment.FileType {
		case "video":
				this.video(attachment)
		}
}

func (this *AttachmentServiceImpl) video(attachment *models.Attachment) bool {
		var err error
		if attachment.Duration == 0 && attachment.Path != "" {
				attachment.Duration, err = libs.GetMp4FileDuration(attachment.GetLocal())
				if err != nil {
						return false
				}
				err = this.model.Update(bson.M{"_id": attachment.Id}, beego.M{"duration": attachment.Duration})
				if err != nil {
						logs.Error(err)
				}
		}
		if "" == this.AutoCoverForVideo(attachment) {
				return false
		}
		return true
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
						attach.Duration = oldAttach.Duration
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

func (this *AttachmentServiceImpl) GetAccessUrl(mediaId string) string {
		var data = this.Get(mediaId)
		if data == nil {
				return ""
		}
		return UrlTicketServiceOf().GetTicketUrlByAttach(data)
}

func (this *AttachmentServiceImpl) AutoCoverForVideo(attachment *models.Attachment,posts...*models.TravelNotes) string {
		if attachment == nil || attachment.FileType != AttachTypeVideo {
				return ""
		}
		if attachment.CoverId != "" {
				return attachment.CoverId.Hex()
		}
		var fs = attachment.GetLocal()
		if !libs.IsExits(fs) {
				return ""
		}
		var (
				ext     = filepath.Ext(fs)
				name    =  fmt.Sprintf("%d.%s", time.Now().Unix(), "jpg")
				storage = strings.Replace(fs, ext,name, 1)
		)
		if plugins.ScreenShot(fs, storage) {
				fd, _ := os.Open(storage)
				defer this.closer(fd)
				data :=beego.M{
						"userId": attachment.UserId,
						"referId": attachment.Id.Hex(), "fileType": AttachTypeImage,
						"filename" : filepath.Base(fd.Name()),
				}
				image := this.Save(fd, data)
				if image == nil {
						return ""
				}
				err := this.UpdateById(attachment.Id.Hex(), beego.M{"coverId": image.Id.Hex()})
				if err != nil {
						logs.Error(err)
				}
				if len(posts)>0 {
						posts[0].Images = append(posts[0].Images,image.Id.Hex())
						return image.Id.Hex()
				}
				if attachment.ReferId != "" && attachment.ReferName == models.TravelNotesTable {
						_ = PostServiceOf().UpdateById(attachment.ReferId, beego.M{"images": []string{image.Id.Hex()}})
				}
				return image.Id.Hex()
		}
		return ""
}

func (this *AttachmentServiceImpl) closer(closer io.Closer) {
		var err = closer.Close()
		if err != nil {
				logs.Error(err)
		}
}
