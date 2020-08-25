package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"time"
)

type CommentModel struct {
		BaseModel
}

// 评论数据
type Comment struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`                                  // ID
		UserId        string        `json:"userId" bson:"userId"`                           // 评论人
		Content       string        `json:"content" bson:"content"`                         // 评论内容
		TargetId      string        `json:"targetId" bson:"targetId"`                       // 评论目标ID
		TargetType    string        `json:"targetType" bson:"targetType"`                   // 评论类型
		Status        int           `json:"status" bson:"status"`                           // 审核状态
		RefersIds     []string      `json:"refersIds,omitempty" bson:"refersIds,omitempty"` // 涉及ID
		ThumbsUpNum   int64         `json:"thumbsUpNum" bson:"thumbsUpNum"`                 // 评论点赞数
		ReviewNum     int64         `json:"reviewNum" bson:"reviewNum"`                     // 评论回复数量
		Sort          int64         `json:"sort" bson:"sort"`                               // 排序
		Tags          []string      `json:"tags" bson:"tags"`                               // 评论标签
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`                     // 评论时间
		UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"`                     // 更新时间
		DeletedAt     int64         `json:"deletedAt" bson:"deletedAt"`                     // 删除时间戳
		dataClassImpl `json:",omitempty" bson:",omitempty"`                                 // 工具类
}

func CommentModelOf() *CommentModel {
		var model = new(CommentModel)
		model._Binder = model
		model.Init()
		return model
}

const (
		CommentTable             = "comments"
		StatusAuditUnPass        = -1
		StatusAuditPass          = 1
		StatusAuditWait          = 0
		StatusOff                = 2
		CommentTargetTypeComment = "post"
		CommentTargetTypeReview  = "comment"
)

var (
		CommentStatusMap = map[int]string{
				StatusAuditWait: "待审核", StatusAuditPass: "审核通过", StatusAuditUnPass: "审核未通过", StatusOff: "下架",
		}
		CommentTypes = []string{
				CommentTargetTypeReview, CommentTargetTypeComment,
		}
)

func NewComment() *Comment {
		return new(Comment)
}

func (this *Comment) CheckType() bool {
		for _, ty := range CommentTypes {
				if this.TargetType == ty {
						return true
				}
		}
		return false
}

func (this *Comment) Load(data map[string]interface{}) *Comment {
		for key, v := range data {
				this.Set(key, v)
		}
		return this
}

func (this *Comment) Set(key string, v interface{}) *Comment {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, v)
		case "userId":
				this.SetString(&this.UserId, v)
		case "content":
				this.SetString(&this.Content, v)
		case "status":
				this.SetNumInt(&this.Status, v)
		case "targetId":
				this.SetString(&this.TargetId, v)
		case "type":
				fallthrough
		case "targetType":
				this.SetString(&this.TargetType, v)
		case "tags":
				this.SetStringArr(&this.Tags, v)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
		case "updatedAt":
				this.SetTime(&this.UpdatedAt, v)
		case "deletedAt":
				this.SetNumIntN(&this.DeletedAt, v)
		case "reviewNum":
				this.SetNumIntN(&this.ReviewNum, v)
		case "sort":
				this.SetNumIntN(&this.Sort, v)
		}
		return this
}

func (this *Comment) GetStatusText() string {
		return CommentStatusMap[this.Status]
}

func (this *Comment) Defaults() *Comment {
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.Sort == 0 {
				this.Sort = 1
		}
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now().Local()
		}
		if this.Tags == nil {
				this.Tags = []string{}
		}
		if this.RefersIds == nil {
				this.RefersIds = []string{}
		}
		return this
}

func (this *Comment) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = beego.M{
				"id":          this.Id.Hex(),
				"sort":        this.Sort,
				"userId":      this.UserId,
				"targetId":    this.TargetId,
				"targetType":  this.TargetType,
				"refersIds":   this.RefersIds,
				"content":     this.Content,
				"tags":        this.Tags,
				"status":      this.Status,
				"statusDesc":  this.GetStatusText(),
				"reviewNum":   this.ReviewNum,
				"thumbsUpNum": this.ThumbsUpNum,
				"updatedAt":   this.UpdatedAt.Unix(),
				"createdAt":   this.CreatedAt.Unix(),
				"deletedAt":   this.DeletedAt,
		}
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}

func (this *Comment) Save() error {
		var (
				id    = this.Id.Hex()
				tmp   = new(User)
				model = CommentModelOf()
				err   = model.GetById(id, tmp)
		)
		if err == nil {
				return model.UpdateById(id, this.M(func(m beego.M) beego.M {
						m = this.excludesKeys(m)
						m["updatedAt"] = time.Now().Local()
						return m
				}))
		}
		return model.Add(this.Defaults())
}

func (this *Comment) excludesKeys(m beego.M) beego.M {
		var excludes = []string{"id", "createdAt", "statusDesc", "updatedAt"}
		for _, k := range excludes {
				delete(m, k)
		}
		return m
}

func (this *CommentModel) TableName() string {
		return CommentTable
}

func (this *CommentModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *CommentModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				this.logs(doc.EnsureIndexKey("userId"))
				this.logs(doc.EnsureIndexKey("postId"))
				this.logs(doc.EnsureIndexKey("commentId"))
				this.logs(doc.EnsureIndexKey("tags"))
				this.logs(doc.EnsureIndexKey("createdAt"))
				this.logs(doc.EnsureIndexKey("status", "sort"))
		}
}
