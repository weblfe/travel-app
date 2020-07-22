package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"time"
)

type PostsModel struct {
		BaseModel
}

func PostsModelOf() *PostsModel {
		var model = new(PostsModel)
		model._Self = model
		model.Init()
		return model
}

// 游记
type TravelNotes struct {
		Id        bson.ObjectId `json:"id" bson:"_id"`
		Title     string        `json:"title" bson:"title"`
		Content   string        `json:"content" bson:"content"`
		Type      int           `json:"type" bson:"type"`
		Images    []string      `json:"images,omitempty" bson:"images,omitempty"`
		UserId    string        `json:"userId" bson:"userId"`
		Videos    []string      `json:"videos,omitempty" bson:"videos,omitempty"`
		Group     string        `json:"group" bson:"group"`
		Tags      []string      `json:"tags" bson:"tags"`
		Status    int           `json:"status" bson:"status"`
		Address   string        `json:"address" bson:"address"`
		Privacy   int           `json:"privacy" bson:"privacy"`
		UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
		CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
		DeletedAt int64         `json:"deletedAt" bson:"deletedAt"`
}

const (
		TravelNotesTable   = "travel_posts"
		PublicPrivacy      = 1
		OnlySelfPrivacy    = 2
		ImageType          = 1
		VideoType          = 2
		ContentType        = 3
		StatusAuditNotPass = -1
		StatusWaitAudit    = 0
		StatusAuditOk      = 1
		StatusAuditOff     = 2
)

var (
		PrivacyMap  = map[int]string{PublicPrivacy: "公开", OnlySelfPrivacy: "仅自己可见"}
		PostTypeMap = map[int]string{ImageType: "图像", VideoType: "视频", ContentType: "文本"}
		StatusMap   = map[int]string{StatusAuditNotPass: "审核不通过", StatusWaitAudit: "待审核", StatusAuditOk: "审核通过", StatusAuditOff: "下架"}
)

func (this *TravelNotes) Load(data bson.M) *TravelNotes {
		for key, v := range data {
				this.Set(key, v)
		}
		return this
}

func (this *TravelNotes) Set(key string, v interface{}) *TravelNotes {
		switch key {
		case "title":
				this.Title = v.(string)
		case "content":
				this.Content = v.(string)
		case "type":
				this.Type = v.(int)
		case "images":
				this.Images = v.([]string)
		case "videos":
				this.Videos = v.([]string)
		case "group":
				this.Group = v.(string)
		case "tags":
				this.Videos = v.([]string)
		case "address":
				this.Address = v.(string)
		case "status":
				this.Status = v.(int)
		case "userId":
				this.UserId = v.(string)
		case "privacy":
				this.Privacy = v.(int)
		case "updatedAt":
				this.UpdatedAt = v.(time.Time)
		case "createdAt":
				this.CreatedAt = v.(time.Time)
		case "deletedAt":
				this.DeletedAt = v.(int64)
		}
		return this
}

func (this *TravelNotes) Defaults() *TravelNotes {
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now().Local()
		}
		if this.Privacy == 0 {
				this.Privacy = PublicPrivacy
		}
		if this.Type == 0 {
				if this.Images != nil && len(this.Images) > 0 {
						this.Type = ImageType
				}
				if this.Videos != nil && len(this.Videos) > 0 {
						this.Type = VideoType
				}
				if this.Type == 0 && this.Content != "" {
						this.Type = ContentType
				}
		}
		return this
}

func (this *TravelNotes) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = beego.M{
				"id":          this.Id.Hex(),
				"title":       this.Title,
				"content":     this.Content,
				"type":        this.Type,
				"typeText":    this.GetType(),
				"images":      this.Images,
				"imagesInfo":  this.GetImages(),
				"userId":      this.UserId,
				"videos":      this.Videos,
				"videosInfo":   this.GetVideos(),
				"tags":        this.Tags,
				"status":      this.Status,
				"statusText":  this.GetState(),
				"group":       this.Group,
				"address":     this.Address,
				"privacy":     this.Privacy,
				"privacyText": this.GetPrivacy(),
				"updatedAt":   this.UpdatedAt.Unix(),
				"createdAt":   this.CreatedAt.Unix(),
				"deletedAt":   this.DeletedAt,
		}
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}

func (this *TravelNotes) GetImages() []*Image {
		if this.Images == nil || len(this.Images) == 0 {
				return nil
		}
		var images []*Image
		for _, id := range this.Images {
				attach := NewAttachment()
				err := AttachmentModelOf().GetById(id, attach)
				if err != nil {
						continue
				}
				images = append(images, attach.Image())
		}
		return images
}

func (this *TravelNotes) GetVideos() []*Video {
		if this.Videos == nil || len(this.Videos) == 0 {
				return nil
		}
		var videos []*Video
		for _, id := range this.Images {
				attach := NewAttachment()
				err := AttachmentModelOf().GetById(id, attach)
				if err != nil {
						continue
				}
				videos = append(videos, attach.Video())
		}
		return videos
}

func (this *TravelNotes) GetPrivacy() string {
		return PrivacyMap[this.Privacy]
}

func (this *TravelNotes) GetType() string {
		return PostTypeMap[this.Type]
}

func (this *TravelNotes) GetState() string {
		return StatusMap[this.Status]
}

func (this *TravelNotes) Save() error {
		var (
				id    = this.Id.Hex()
				tmp   = new(TravelNotes)
				model = PostsModelOf()
				err   = model.GetById(id, tmp)
		)
		if err != nil {
				return model.UpdateById(id, this.M(func(m beego.M) beego.M {
						delete(m, "id")
						delete(m, "createdAt")
						m["updatedAt"] = time.Now().Local()
						return m
				}))
		}
		return model.Add(this)
}

func (this *TravelNotes) IsEmpty() bool  {
		if this.Content == "" || (len(this.Videos) == 0  && len(this.Images)==0) {
				return true
		}
		return false
}

func (this *PostsModel) TableName() string {
		return TravelNotesTable
}

func (this *PostsModel) CreateIndex() {
		_ = this.Collection().EnsureIndexKey("title")
		_ = this.Collection().EnsureIndexKey("type")
		_ = this.Collection().EnsureIndexKey("userId")
		_ = this.Collection().EnsureIndexKey("tags")
		_ = this.Collection().EnsureIndexKey("group")
		_ = this.Collection().EnsureIndexKey("address", "privacy")
}
