package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"regexp"
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
		Id          bson.ObjectId `json:"id" bson:"_id"`
		Title       string        `json:"title" bson:"title"`                       // 标题
		Content     string        `json:"content" bson:"content"`                   // 内容
		Type        int           `json:"type" bson:"type"`                         // 类型
		Images      []string      `json:"images,omitempty" bson:"images,omitempty"` // 图片ID
		UserId      string        `json:"userId" bson:"userId"`                     // 用户ID
		Videos      []string      `json:"videos,omitempty" bson:"videos,omitempty"` // 视频ID
		Group       string        `json:"group" bson:"group"`                       // 分组类型名
		Tags        []string      `json:"tags" bson:"tags"`                         // 标签ID
		Status      int           `json:"status" bson:"status"`                     // 审核状态
		Address     string        `json:"address" bson:"address"`                   // 地址
		Privacy     int           `json:"privacy" bson:"privacy"`                   // 是否公开
		ThumbsUpNum int64         `json:"thumbsUpNum" bson:"thumbsUpNum"`           // 点赞数
		CommentNum  int64         `json:"commentNum" bson:"commentNum"`             // 评论数
		Score       int64         `json:"score" bson:"score"`                       // 作品评分
		UpdatedAt   time.Time     `json:"updatedAt" bson:"updatedAt"`               // 更新时间
		CreatedAt   time.Time     `json:"createdAt" bson:"createdAt"`               // 创建时间
		DeletedAt   int64         `json:"deletedAt" bson:"deletedAt"`               // 删除时间
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
				"images":      this.getImages(),
				"imagesInfo":  this.GetImages(),
				"userId":      this.UserId,
				"videos":      this.getVideos(),
				"videosInfo":  this.GetVideos(),
				"tags":        this.getTags(),
				"tagsText":    this.GetTagsText(),
				"status":      this.Status,
				"statusText":  this.GetState(),
				"group":       this.Group,
				"address":     this.Address,
				"privacy":     this.Privacy,
				"commentNum":  this.CommentNum,
				"thumbsUpNum": this.ThumbsUpNum,
				"score":       this.Score,
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

// 移除多余无需更新字段
func (this *TravelNotes) removeUpdateExcludes(m beego.M) beego.M {
		var keys = []string{"imagesInfo", "videosInfo", "tagsText", "statusText", "privacyText", "id", "createdAt"}
		for _, k := range keys {
				delete(m, k)
		}
		return m
}

// 获取标签描述
func (this *TravelNotes) GetTagsText() []string {
		if this.Tags == nil || len(this.Tags) == 0 {
				return []string{}
		}
		var (
				err    error
				result = make([]string, 1)
				arr    = make([]bson.ObjectId, 0)
				tags   = make([]*Tag, len(this.Tags))
				regex  = regexp.MustCompile(`^\w+$`)
		)
		tags = tags[:0]
		for _, tag := range this.Tags {
				if regex.MatchString(tag) {
						arr = append(arr, bson.ObjectIdHex(tag))
				}
		}
		err = TagsModelOf().Gets(bson.M{"_id": beego.M{"$in": arr}}, &tags)
		if err != nil {
				return []string{}
		}
		result = result[:0]
		for _, it := range tags {
				result = append(result, it.Name)
		}
		return result
}

func (this *TravelNotes) getTags() []string {
		if this.Tags == nil || len(this.Tags) <= 0 {
				return []string{}
		}
		var (
				arr   = make([]string, 2)
				regex = regexp.MustCompile(`^\w+$`)
		)
		arr = arr[:0]
		for _, tag := range this.Tags {
				if regex.MatchString(tag) {
						arr = append(arr, tag)
				}
		}
		return arr
}

func (this *TravelNotes) getVideos() []string {
		if this.Videos == nil {
				return []string{}
		}
		return this.Videos
}

func (this *TravelNotes) getImages() []string {
		if this.Images == nil {
				return []string{}
		}
		return this.Images
}

func (this *TravelNotes) GetImages() []*Image {
		if this.Images == nil || len(this.Images) == 0 {
				return []*Image{}
		}
		var (
				images      []*Image
				ids         []bson.ObjectId
				attachArr   = make([]*Attachment, 2)
				attachModel = AttachmentModelOf()
		)
		for _, v := range this.Images {
				ids = append(ids, bson.ObjectIdHex(v))
		}
		attachArr = attachArr[:0]
		var err = attachModel.Gets(bson.M{"_id": bson.M{"$in": ids}}, &attachArr)
		if err == nil {
				for _, attach := range attachArr {
						if attach == nil {
								continue
						}
						images = append(images, attach.Image())
				}
		}
		if images == nil {
				return []*Image{}
		}
		return images
}

func (this *TravelNotes) GetVideos() []*Video {
		if this.Videos == nil || len(this.Videos) == 0 {
				return []*Video{}
		}
		var (
				videos      []*Video
				ids         []bson.ObjectId
				attachArr   = make([]*Attachment, 2)
				attachModel = AttachmentModelOf()
		)
		for _, v := range this.Videos {
				ids = append(ids, bson.ObjectIdHex(v))
		}
		attachArr = attachArr[:0]
		var err = attachModel.Gets(bson.M{"_id": bson.M{"$in": ids}}, &attachArr)
		if err == nil {
				for _, attach := range attachArr {
						if attach == nil {
								continue
						}
						videos = append(videos, attach.Video())
				}
		}
		if videos == nil {
				return []*Video{}
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
		if err == nil {
				return model.UpdateById(id, this.M(func(m beego.M) beego.M {
						m = this.removeUpdateExcludes(m)
						m["updatedAt"] = time.Now().Local()
						return m
				}))
		}
		return model.Add(this)
}

func (this *TravelNotes) IsEmpty() bool {
		if this.Content == "" || (len(this.Videos) == 0 && len(this.Images) == 0) {
				return true
		}
		return false
}

func (this *PostsModel) TableName() string {
		return TravelNotesTable
}

func (this *PostsModel) CreateIndex() {
		//	_ = this.Collection().EnsureIndexKey("title")
		_ = this.Collection().EnsureIndexKey("type")
		_ = this.Collection().EnsureIndexKey("userId")
		_ = this.Collection().EnsureIndexKey("tags")
		_ = this.Collection().EnsureIndexKey("group")
		_ = this.Collection().EnsureIndexKey("address", "privacy")

		_ = this.Collection().EnsureIndexKey("thumbsUpNum", "commentNum", "score")
		// null unique username
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:              []string{"$text:content"},
				DefaultLanguage:  "chinese",
				LanguageOverride: "language",
		})
}

// 增加
func (this *PostsModel) Incr(id string, typ string, num ...int) error {
		if len(num) == 0 {
				num = append(num, 1)
		}
		return this.IncrBy(bson.M{"_id": bson.ObjectIdHex(id)}, beego.M{typ: int64(num[0])})
}
