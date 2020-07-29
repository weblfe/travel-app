package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"strings"
		"time"
)

type CommentModel struct {
		BaseModel
}

// 评论数据
type Comment struct {
		Id          bson.ObjectId `json:"id" bson:"_id"`                                  // ID
		UserId      string        `json:"userId" bson:"userId"`                           // 评论人
		Content     string        `json:"content" bson:"content"`                         // 评论内容
		PostId      bson.ObjectId `json:"postId" bson:"postId"`                           // 评论文章ID
		Status      int           `json:"status" bson:"status"`                           // 审核状态
		CommentId   bson.ObjectId `json:"commentId,omitempty" bson:"commentId,omitempty"` // 评论 评论的ID
		ThumbsUpNum int64         `json:"thumbsUpNum" bson:"thumbsUpNum"`                 // 评论点赞数
		Sort        int64         `json:"sort" bson:"sort"`                               // 排序
		Tags        []string      `json:"tags" bson:"tags"`                               // 评论标签
		CreatedAt   time.Time     `json:"createdAt" bson:"createdAt"`                     // 评论时间
		UpdatedAt   time.Time     `json:"updatedAt" bson:"updatedAt"`                     // 更新时间
		DeletedAt   int64         `json:"deletedAt" bson:"deletedAt"`                     // 删除时间戳
}

func CommentModelOf() *CommentModel {
		var model = new(CommentModel)
		model._Self = model
		model.Init()
		return model
}

const (
		CommentTable      = "comments"
		StatusAuditPass   = 1
		StatusAuditWait   = 0
		StatusAuditUnPass = 2
		StatusOff         = -1
)

var (
		CommentStatusMap = map[int]string{
				StatusAuditWait: "待审核", StatusAuditPass: "审核通过", StatusAuditUnPass: "审核未通过", StatusOff: "下架",
		}
)

func (this *Comment) Load(data map[string]interface{}) *Comment {
		for key, v := range data {
				this.Set(key, v)
		}
		return this
}

func (this *Comment) Set(key string, v interface{}) *Comment {
		switch key {
		case "id":
				if str, ok := v.(string); ok && str != "" {
						this.Id = bson.ObjectIdHex(str)
				}
				if obj, ok := v.(bson.ObjectId); ok && obj != "" {
						this.Id = obj
				}
		case "userId":
				this.UserId = v.(string)
		case "content":
				this.Content = v.(string)
		case "status":
				this.Status = v.(int)
		case "commentId":
				if str, ok := v.(string); ok && str != "" {
						this.CommentId = bson.ObjectIdHex(str)
				}
				if obj, ok := v.(bson.ObjectId); ok && obj != "" {
						this.CommentId = obj
				}
		case "postId":
				if str, ok := v.(string); ok && str != "" {
						this.PostId = bson.ObjectIdHex(str)
				}
				if obj, ok := v.(bson.ObjectId); ok && obj != "" {
						this.PostId = obj
				}
		case "tags":
				if str, ok := v.(string); ok && str != "" {
						this.Tags = strings.SplitN(str, ",", -1)
				}
				if arr, ok := v.([]string); ok && len(arr) > 0 {
						this.Tags = arr
				}
		case "createdAt":
				this.CreatedAt = v.(time.Time)
		case "updatedAt":
				this.UpdatedAt = v.(time.Time)
		case "deletedAt":
				this.DeletedAt = v.(int64)
		case "sort":
				this.Sort = v.(int64)
		}
		return this
}

func (this *Comment) GetStatus() string {
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
		return this
}

func (this *Comment) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = beego.M{
				"id":          this.Id.Hex(),
				"sort":        this.Sort,
				"userId":      this.UserId,
				"postId":      this.PostId.Hex(),
				"comment":     this.CommentId.Hex(),
				"content":     this.Content,
				"tags":        this.Tags,
				"status":      this.Status,
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
		if err != nil {
				return model.UpdateById(id, this.M(func(m beego.M) beego.M {
						delete(m, "id")
						delete(m, "createdAt")
						m["updatedAt"] = time.Now().Local()
						return m
				}))
		}
		return model.Add(this.Defaults())
}

func (this *CommentModel) TableName() string {
		return CommentTable
}

func (this *CommentModel) CreateIndex() {
		_ = this.Collection().EnsureIndexKey("userId")
		_ = this.Collection().EnsureIndexKey("postId")
		_ = this.Collection().EnsureIndexKey("commentId")
		_ = this.Collection().EnsureIndexKey("tags")
		_ = this.Collection().EnsureIndexKey("createdAt")
		_ = this.Collection().EnsureIndexKey("status", "sort")
}
