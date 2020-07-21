package models

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/config/env"
		"github.com/astaxie/beego/orm"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"path/filepath"
		"strings"
		"time"
)

type AttachmentModel struct {
		BaseModel
}

// 附件模型
type Attachment struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`                                          // id media_id
		FileName      string        `json:"filename" bson:"filename"`                               // 文件名
		Hash          string        `json:"hash" bson:"hash"`                                       // 文件hash值
		Ticket        string        `json:"ticket" bson:"ticket"`                                   // 文件上传时的密钥
		AppId         string        `json:"appId" bson:"appId"`                                     // 文件上传的应用
		UserId        bson.ObjectId `json:"userId,omitempty" bson:"userId,omitempty"`               // 文件上传用户
		ExtrasInfo    bson.M        `json:"extrasInfo,omitempty" bson:"extrasInfo,omitempty"`       // 文件扩展信息
		Tags          []string      `json:"tags,omitempty" bson:"tags,omitempty"`                   // 文件标签
		Url           string        `json:"url" bson:"url"`                                         // 远程访问链接
		Cdn           string        `json:"cdn,omitempty" bson:"cdn,omitempty"`                     // cdn 服务名
		CdnUrl        string        `json:"cdnUrl,omitempty" bson:"cdnUrl,omitempty"`               // cdn访问链接
		OssBucket     string        `json:"ossBucket,omitempty" bson:"ossBucket,omitempty"`         // oss bucket
		Oss           string        `json:"oss,omitempty" bson:"oss,omitempty"`                     // oss 服务名
		AccessTimes   int64         `json:"accessTimes,omitempty" bson:"accessTimes,omitempty"`     // 被访问次数
		DownloadTimes int64         `json:"downloadTimes,omitempty" bson:"downloadTimes,omitempty"` // 被下载次数
		Path          string        `json:"path" bson:"path"`                                       // 系统本地存储路径
		ReferName     string        `json:"referName" bson:"referName"`                             // 记录涉及的document名
		ReferId       string        `json:"referId" bson:"referId"`                                 // 记录涉及的document的ID
		Size          int64         `json:"size" bson:"size"`                                       // 文件大小 单位: byte
		SizeText      string        `json:"sizeText" bson:"sizeText"`                               // 带单的文件大小 eg: ..1G,120MB,1KB,1B,1byte
		FileType      string        `json:"fileType" bson:"fileType"`                               // 文件类型 [doc,image,mp4,mp3,txt....]
		Status        int           `json:"status" bson:"status"`                                   // 文件状态
		Privately     bool          `json:"privately" bson:"privately"`                             // 文件是否私有
		Watermark     bool          `json:"watermark" bson:"watermark"`                             // 文件是否有水印
		UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"`                             // 记录更新时间
		Duration      time.Duration `json:"duration,omitempty" bson:"duration,omitempty"`           // 音视频文件时长
		CoverId       bson.ObjectId `json:"coverId,omitempty" bson:"coverId,omitempty"`             // 音视频文件封面
		Width         int           `json:"width,omitempty" bson:"width,omitempty"`                 // 图片文件时宽
		Height        int           `json:"height,omitempty" bson:"height,omitempty"`               // 图片文件时高
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`                             // 创建时间
		DeletedAt     int64         `json:"deletedAt" bson:"deletedAt"`                             // 删除时间
}

// 图片
type Image struct {
		MediaId  string `json:"mediaId" bson:"mediaId"`
		Url      string `json:"url" bson:"url"`
		Size     int    `json:"size" bson:"size"`
		SizeText string `json:"sizeText" bson:"sizeText"`
		Width    int    `json:"width" bson:"width"`
		Height   int    `json:"height" bson:"height"`
}

// 视频
type Video struct {
		MediaId  string        `json:"mediaId" bson:"mediaId"`
		Url      string        `json:"url" bson:"url"`
		Size     int           `json:"size" bson:"size"`
		SizeText string        `json:"sizeText" bson:"sizeText"`
		Duration time.Duration `json:"width" bson:"width"`
		CoverId  string        `json:"coverId" bson:"coverId"`
}

const (
		AttachmentTable = "attachments"
)

func NewAttachment() *Attachment {
		return new(Attachment)
}

func AttachmentModelOf() *AttachmentModel {
		var model = new(AttachmentModel)
		model._Self = model
		model.Init()
		return model
}

func (this *Attachment) Load(data map[string]interface{}) *Attachment {
		for key, v := range data {
				this.set(key, v)
		}
		return this
}

func (this *Attachment) set(key string, v interface{}) *Attachment {
		switch key {
		case "Filename":
				fallthrough
		case "filename":
				this.FileName = v.(string)
		case "Hash":
				fallthrough
		case "hash":
				this.Hash = v.(string)
		case "Ticket":
				fallthrough
		case "ticket":
				this.Ticket = v.(string)
		case "appId":
				fallthrough
		case "AppId":
				this.AppId = v.(string)
		case "userId":
				fallthrough
		case "UserId":
				if str, ok := v.(string); ok && str != "" {
						this.UserId = bson.ObjectIdHex(str)
				}
				if id, ok := v.(bson.ObjectId); ok {
						this.UserId = id
				}
		case "coverId":
				fallthrough
		case "CoverId":
				if str, ok := v.(string); ok && str != "" {
						this.CoverId = bson.ObjectIdHex(str)
				}
				if id, ok := v.(bson.ObjectId); ok {
						this.CoverId = id
				}
		case "ExtrasInfo":
				fallthrough
		case "extrasInfo":
				if m, ok := v.(map[string]interface{}); ok {
						this.ExtrasInfo = bson.M(m)
				}
				if m, ok := v.(beego.M); ok {
						this.ExtrasInfo = bson.M(m)
				}
				if m, ok := v.(bson.M); ok {
						this.ExtrasInfo = m
				}
		case "Tags":
				fallthrough
		case "tags":
				if str, ok := v.(string); ok {
						this.Tags = strings.SplitN(str, ",", -1)
				}
				if strArr, ok := v.([]string); ok {
						this.Tags = strArr
				}
		case "Url":
				fallthrough
		case "url":
				this.Url = v.(string)
		case "Cdn":
				fallthrough
		case "cdn":
				this.Cdn = v.(string)
		case "CdnUrl":
				fallthrough
		case "cdnUrl":
				this.CdnUrl = v.(string)
		case "OssBucket":
				fallthrough
		case "ossBucket":
				this.OssBucket = v.(string)
		case "Oss":
				fallthrough
		case "oss":
				this.Oss = v.(string)
		case "accessTimes":
				fallthrough
		case "AccessTimes":
				this.AccessTimes = orm.ToInt64(v)
		case "downloadTimes":
				fallthrough
		case "DownloadTimes":
				this.DownloadTimes = orm.ToInt64(v)
		case "path":
				fallthrough
		case "Path":
				this.Path = orm.ToStr(v)
		case "referName":
		case "ReferName":
				this.ReferName = orm.ToStr(v)
		case "ReferId":
				fallthrough
		case "referId":
				this.ReferId = orm.ToStr(v)
		case "Size":
				fallthrough
		case "size":
				this.Size = orm.ToInt64(v)
		case "sizeText":
				fallthrough
		case "SizeText":
				this.SizeText = orm.ToStr(v)
		case "fileType":
				fallthrough
		case "FileType":

				this.FileType = orm.ToStr(v)
		case "status":
				fallthrough
		case "Status":
				this.Status = int(orm.ToInt64(v))
		case "privately":
				fallthrough
		case "Privately":
				if n, ok := v.(int); ok {
						if n > 0 {
								this.Privately = true
						}
				}
				if n, ok := v.(string); ok {
						if n == "true" || n == "True" || n == "TRUE" {
								this.Privately = true
						}
				}
				if b, ok := v.(bool); ok {
						this.Privately = b
				}
		case "watermark":
				fallthrough
		case "Watermark":
				if n, ok := v.(int); ok {
						if n > 0 {
								this.Watermark = true
						}
				}
				if n, ok := v.(string); ok {
						if n == "true" || n == "True" || n == "TRUE" {
								this.Watermark = true
						}
				}
				if b, ok := v.(bool); ok {
						this.Watermark = b
				}
		case "duration":
				fallthrough
		case "Duration":
				if n, ok := v.(int64); ok {
						this.Duration = time.Duration(n)
				}
				if n, ok := v.(time.Duration); ok {
						this.Duration = n
				}
		case "width":
				fallthrough
		case "Width":
				this.Width = int(orm.ToInt64(v))
		case "height":
				fallthrough
		case "Height":
				this.Height = int(orm.ToInt64(v))
		case "createdAt":
				fallthrough
		case "CreatedAt":
				this.CreatedAt = v.(time.Time)
		case "deletedAt":
				fallthrough
		case "DeletedAt":
				this.DeletedAt = orm.ToInt64(v)
		}
		return this
}

func (this *Attachment) Defaults() *Attachment {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.UserId == "" {
				this.UserId = ""
		}
		if this.Url == "" {
				if this.CdnUrl == "" && this.Path != "" && this.FileName != "" {
						this.Url = fmt.Sprintf(
								"%s://%s/%s/%s",
								env.Get("SERVER_SCHEMA", "http"),
								env.Get("SERVER_DOMAIN", "localhost"),
								env.Get("ATTACHMENT_PATH", "attachments"),
								this.Id.Hex(),
						)
				}
		}
		if this.SizeText == "" && this.Size != 0 {
				this.SizeText = libs.FormatFileSize(this.Size)
		}
		if this.ExtrasInfo == nil {
				this.ExtrasInfo = make(bson.M)
		}
		if this.FileType == "" && this.FileName != "" {
				this.FileType = libs.GetFileType(this.FileName)
		}
		if _, ok := this.ExtrasInfo["extension"]; !ok && this.FileName != "" {
				this.ExtrasInfo["extension"] = filepath.Ext(this.FileName)[1:]
		}
		return this
}

func (this *Attachment) M(filters ...func(m beego.M) beego.M) beego.M {
		var m = beego.M{
				"id":            this.Id.Hex(),
				"filename":      this.FileName,
				"hash":          this.Hash,
				"ticket":        this.Ticket,
				"appId":         this.AppId,
				"userId":        this.UserId.Hex(),
				"extrasInfo":    this.ExtrasInfo,
				"tags":          this.Tags,
				"url":           this.Url,
				"cdn":           this.Cdn,
				"cdnUrl":        this.CdnUrl,
				"ossBucket":     this.OssBucket,
				"oss":           this.Oss,
				"accessTimes":   this.AccessTimes,
				"downloadTimes": this.DownloadTimes,
				"path":          this.Path,
				"referName":     this.ReferName,
				"referId":       this.ReferId,
				"size":          this.Size,
				"sizeText":      this.SizeText,
				"fileType":      this.FileType,
				"status":        this.Status,
				"privately":     this.Privately,
				"watermark":     this.Watermark,
				"duration":      this.Duration,
				"width":         this.Width,
				"height":        this.Height,
				"coverId":       this.CoverId,
				"updatedAt":     this.UpdatedAt.Unix(),
				"createdAt":     this.CreatedAt.Unix(),
				"deletedAt":     this.DeletedAt,
		}
		if len(filters) != 0 {
				for _, filter := range filters {
						if len(m) == 0 {
								return m
						}
						m = filter(m)
				}
		}
		return m
}

func (this *Attachment) GetUrl() string {
		if this.CdnUrl != "" {
				return this.CdnUrl
		}
		return this.Url
}

func (this *Attachment) Image() *Image {
		var image = new(Image)
		image.Height = this.Height
		image.MediaId = this.Id.Hex()
		image.Width = this.Width
		image.SizeText = this.SizeText
		image.Size = int(this.Size)
		image.Url = this.GetUrl()
		return image
}

func (this *Attachment) Video() *Video {
		var video = new(Video)
		video.Duration = this.Duration
		if this.CoverId != "" {
				video.CoverId = this.CoverId.Hex()
		}
		video.MediaId = this.Id.Hex()
		video.SizeText = this.SizeText
		video.Size = int(this.Size)
		video.Url = this.GetUrl()
		return video
}

func (this *AttachmentModel) CreateIndex() {
		_ = this.Collection().EnsureIndexKey("appId")
		_ = this.Collection().EnsureIndexKey("filename")
		_ = this.Collection().EnsureIndexKey("referName", "referId")
		_ = this.Collection().EnsureIndexKey("size")
		_ = this.Collection().EnsureIndexKey("fileType")
		_ = this.Collection().EnsureIndexKey("status")
		_ = this.Collection().EnsureIndexKey("deletedAt")
		// unique mobile
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"access_count", "downloadTimes"},
				Unique: false,
				Sparse: true,
		})
}

func (this *AttachmentModel) TableName() string {
		return AttachmentTable
}

func (this *Video) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = beego.M{
				"url":      this.Url,
				"size":     this.Size,
				"mediaId":  this.MediaId,
				"sizeText": this.SizeText,
				"coverId":  this.CoverId,
				"duration": this.Duration,
		}
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}

func (this *Image) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = beego.M{
				"mediaId":  this.MediaId,
				"url":      this.Url,
				"size":     this.Size,
				"sizeText": this.SizeText,
				"width":    this.Width,
				"height":   this.Height,
		}
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}

func (this *AttachmentModel) GetByMediaId(id string) (*Attachment, error) {
		var att = NewAttachment()
		if id == "" {
				return nil, common.NewErrors("empty id")
		}
		err := this.GetById(id, att)
		if err == nil {
				return att, nil
		}
		return nil, err
}
