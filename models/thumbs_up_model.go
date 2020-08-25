package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"time"
)

type ThumbsUpModel struct {
		BaseModel
}

// 点赞数据
type ThumbsUp struct {
		Id        bson.ObjectId `json:"id" bson:"_id"`
		Status    int           `json:"status" bson:"status"`
		Type      string        `json:"type" bson:"type"` // 类型
		UserId    string        `json:"userId" bson:"userId"`
		TypeId    string        `json:"typeId" bson:"typeId"`
		Count     int           `json:"count" bson:"count"`
		CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
		UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

const (
		ThumbsTypePost    = "post"    // 游记点
		ThumbsTypeComment = "comment" // 评论点赞
)

func (this *ThumbsUp) Load(data map[string]interface{}) *ThumbsUp {
		for key, v := range data {
				this.Set(key, v)
		}
		return this
}

func (this *ThumbsUp) Set(key string, v interface{}) *ThumbsUp {
		switch key {
		case "id":
				if str, ok := v.(string); ok && str != "" {
						this.Id = bson.ObjectIdHex(str)
				}
				if obj, ok := v.(bson.ObjectId); ok && obj != "" {
						this.Id = obj
				}
		case "status":
				this.Status = v.(int)
		case "type":
				this.Type = v.(string)
		case "userId":
				this.UserId = v.(string)
		case "typeId":
				this.TypeId = v.(string)
		case "count":
				this.Count = v.(int)
		case "createdAt":
				this.CreatedAt = v.(time.Time)
		case "updatedAt":
				this.UpdatedAt = v.(time.Time)
		}
		return this
}

func (this *ThumbsUp) Defaults() *ThumbsUp {
		if this.Status == 0 && this.Id == "" && this.CreatedAt.IsZero() {
				this.Status = 1
		}
		if this.Count == 0 && this.Id == "" && this.CreatedAt.IsZero() {
				this.Count = 1
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now().Local()
		}
		return this
}

func (this *ThumbsUp) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = beego.M{
				"id":        this.Id.Hex(),
				"status":    this.Status,
				"count":     this.Count,
				"type":      this.Type,
				"typeId":    this.TypeId,
				"userId":    this.UserId,
				"createdAt": this.CreatedAt,
				"updatedAt": this.UpdatedAt,
		}
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}

func (this *ThumbsUp) Save() error {
		var (
				err   error
				id    = this.Id.Hex()
				tmp   = new(ThumbsUp)
				model = ThumbsUpModelOf()
		)
		if id != "" {
				err = model.GetById(id, tmp)
		}
		if err == nil {
				return model.UpdateById(id, this.M(func(m beego.M) beego.M {
						delete(m, "id")
						delete(m, "createdAt")
						m["updatedAt"] = time.Now().Local()
						return m
				}))
		}
		return model.Add(this.Defaults())
}

func ThumbsUpModelOf() *ThumbsUpModel {
		var model = new(ThumbsUpModel)
		model._Binder = model
		model.Init()
		return model
}

const (
		ThumbsUpTable = "thumbs_up"
)

func (this *ThumbsUpModel) TableName() string {
		return ThumbsUpTable
}

func (this *ThumbsUpModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *ThumbsUpModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"type", "typeId", "userId"},
						Unique: true,
						Sparse: true,
				}))
				this.logs(doc.EnsureIndexKey("state"))
				this.logs(doc.EnsureIndexKey("gender"))
				this.logs(doc.EnsureIndexKey("nickname"))
		}
}
