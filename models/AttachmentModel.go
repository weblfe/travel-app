package models

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/orm"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"strings"
		"time"
)

type AttachmentModel struct {
		BaseModel
}

// 附件模型
type Attachment struct {
		Id            bson.ObjectId `json:"id" bson:"id"`                                             // id media_id
		FileName      string        `json:"filename" bson:"filename"`                                 // 文件名
		Hash          string        `json:"hash" bson:"hash"`                                         // 文件hash值
		Ticket        string        `json:"ticket" bson:"ticket"`                                     // 文件上传时的密钥
		AppId         string        `json:"app_id" bson:"app_id"`                                     // 文件上传的应用
		UserId        bson.ObjectId `json:"user_id,omitempty" bson:"user_id,omitempty"`               // 文件上传用户
		ExtrasInfo    bson.M        `json:"extras_info,omitempty" bson:"extras_info,omitempty"`       // 文件扩展信息
		Tags          []string      `json:"tags,omitempty" bson:"tags,omitempty"`                     // 文件标签
		Url           string        `json:"url" bson:"url"`                                           // 远程访问链接
		Cdn           string        `json:"cdn,omitempty" bson:"cdn,omitempty"`                       // cdn 服务名
		CdnUrl        string        `json:"cdn_url,omitempty" bson:"cdn_url,omitempty"`               // cdn访问链接
		OssBucket     string        `json:"oss_bucket,omitempty" bson:"oss_bucket,omitempty"`         // oss bucket
		Oss           string        `json:"oss,omitempty" bson:"oss,omitempty"`                       // oss 服务名
		AccessTimes   int64         `json:"access_times,omitempty" bson:"access_times,omitempty"`     // 被访问次数
		DownloadTimes int64         `json:"download_times,omitempty" bson:"download_times,omitempty"` // 被下载次数
		Path          string        `json:"path" bson:"path"`                                         // 系统本地存储路径
		ReferName     string        `json:"refer_name" bson:"refer_name"`                             // 记录涉及的document名
		ReferId       string        `json:"refer_id" bson:"refer_id"`                                 // 记录涉及的document的ID
		Size          int64         `json:"size" bson:"size"`                                         // 文件大小 单位: byte
		SizeText      string        `json:"size_text" bson:"size_text"`                               // 带单的文件大小 eg: ..1G,120MB,1KB,1B,1byte
		FileType      string        `json:"file_type" bson:"file_type"`                               // 文件类型 [doc,image,mp4,mp3,txt....]
		Status        int           `json:"status" bson:"status"`                                     // 文件状态
		Privately     bool          `json:"privately" bson:"privately"`                               // 文件是否私有
		Watermark     bool          `json:"watermark" bson:"watermark"`                               // 文件是否有水印
		UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`                             // 记录更新时间
		Duration      time.Duration `json:"duration,omitempty" bson:"duration,omitempty"`             // 音视频文件时长
		Width         int           `json:"width,omitempty" bson:"width,omitempty"`                   // 图片文件时宽
		Height        int           `json:"height,omitempty" bson:"height,omitempty"`                 // 图片文件时高
		CreatedAt     time.Time     `json:"created_at" bson:"created_at"`                             // 创建时间
		DeletedAt     int64         `json:"deleted_at" bson:"deleted_at"`                             // 删除时间
}

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
		case "app_id":
				fallthrough
		case "AppId":
				this.AppId = v.(string)
		case "user_id":
				fallthrough
		case "UserId":
				if str, ok := v.(string); ok {
						this.UserId = bson.ObjectIdHex(str)
				}
				if id, ok := v.(bson.ObjectId); ok {
						this.UserId = id
				}
		case "ExtrasInfo":
				fallthrough
		case "extras_info":
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
		case "cdn_url":
				this.CdnUrl = v.(string)
		case "OssBucket":
				fallthrough
		case "oss_bucket":
				this.OssBucket = v.(string)
		case "Oss":
				fallthrough
		case "oss":
				this.Oss = v.(string)
		case "access_times":
				fallthrough
		case "AccessTimes":
				this.AccessTimes = orm.ToInt64(v)
		case "download_times":
				fallthrough
		case "DownloadTimes":
				this.DownloadTimes = orm.ToInt64(v)
		case "path":
				fallthrough
		case "Path":
				this.Path = orm.ToStr(v)
		case "refer_name":
		case "ReferName":
				this.ReferName = orm.ToStr(v)
		case "ReferId":
				fallthrough
		case "refer_id":
				this.ReferId = orm.ToStr(v)
		case "Size":
				fallthrough
		case "size":
				this.Size = orm.ToInt64(v)
		case "size_text":
				fallthrough
		case "SizeText":
				this.SizeText = orm.ToStr(v)
		case "file_type":
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
		case "created_at":
				fallthrough
		case "CreatedAt":
				this.CreatedAt = v.(time.Time)
		case "deleted_at":
				fallthrough
		case "DeletedAt":
				this.DeletedAt = orm.ToInt64(v)
		}
		return this
}

func (this *Attachment) Defaults() *Attachment {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now()
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		return this
}

func (this *Attachment) M(filters ...func(m beego.M) beego.M) beego.M {
		var m = beego.M{
				"id":             this.Id.Hex(),
				"filename":       this.FileName,
				"hash":           this.Hash,
				"ticket":         this.Ticket,
				"app_id":         this.AppId,
				"user_id":        this.UserId.Hex(),
				"extras_info":    this.ExtrasInfo,
				"tags":           this.Tags,
				"url":            this.Url,
				"cdn":            this.Cdn,
				"cdn_url":        this.CdnUrl,
				"oss_bucket":     this.OssBucket,
				"oss":            this.Oss,
				"access_times":   this.AccessTimes,
				"download_times": this.DownloadTimes,
				"path":           this.Path,
				"refer_name":     this.ReferName,
				"refer_id":       this.ReferId,
				"size":           this.Size,
				"size_text":      this.SizeText,
				"file_type":      this.FileType,
				"status":         this.Status,
				"privately":      this.Privately,
				"watermark":      this.Watermark,
				"duration":       this.Duration,
				"width":          this.Width,
				"height":         this.Height,
				"created_at":     this.CreatedAt,
				"deleted_at":     this.DeletedAt,
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

func (this *AttachmentModel) CreateIndex() {
		_ = this.Collection().EnsureIndexKey("app_id")
		_ = this.Collection().EnsureIndexKey("filename")
		_ = this.Collection().EnsureIndexKey("refer_name", "refer_id")
		_ = this.Collection().EnsureIndexKey("size")
		_ = this.Collection().EnsureIndexKey("file_type")
		_ = this.Collection().EnsureIndexKey("status")
		_ = this.Collection().EnsureIndexKey("deleted_at")
		// unique mobile
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"access_count", "download_times"},
				Unique: false,
				Sparse: true,
		})
}

func (this *AttachmentModel) GetByMediaId(id string) (*Attachment, error) {
		var att = NewAttachment()
		err := this.GetById(id, att)
		if err == nil {
				return att, nil
		}
		return nil, err
}
