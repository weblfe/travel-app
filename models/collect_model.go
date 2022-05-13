package models

import (
	"github.com/astaxie/beego"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/weblfe/travel-app/transforms"
	"log"
	"time"
)

// Collect  收藏记录
type Collect struct {
	Id            bson.ObjectId                          `json:"id" bson:"_id"`                // ID
	UserId        string                                 `json:"userId" bson:"userId"`         // 用户ID
	TargetId      string                                 `json:"targetId" bson:"targetId"`     // 收藏目标ID
	TargetType    string                                 `json:"targetType" bson:"targetType"` // 收藏类型
	Status        int                                    `json:"status" bson:"status"`         // 状态
	Versions      []string                               `json:"versions" bson:"versions"`     // 状态版本变化历史
	CreatedAt     time.Time                              `json:"createdAt" bson:"createdAt"`   // 创建时间
	UpdatedAt     time.Time                              `json:"updatedAt" bson:"updatedAt"`   // 更新时间
	DeletedAt     int64                                  `json:"deletedAt" bson:"deletedAt"`   // 删除时间
	dataClassImpl `bson:",omitempty"  json:",omitempty"` // 工具类
}

// CollectModel 收藏记录model
type CollectModel struct {
	BaseModel
}

const (
	CollectTable          = "user_collects"
	CollectTargetTypePost = "post"
	CollectTargetVideo    = "video"
	StatusActivity        = 1
)

func NewCollect() *Collect {
	var data = new(Collect)
	data.Init()
	return data
}

func CollectModelOf() *CollectModel {
	var model = new(CollectModel)
	return model.init()
}

func (this *Collect) data() beego.M {
	return beego.M{
		"id":         this.Id.Hex(),
		"userId":     this.UserId,
		"targetId":   this.TargetId,
		"targetType": this.TargetType,
		"status":     this.Status,
		"versions":   this.getVersions(),
		"createdAt":  this.CreatedAt.Unix(),
		"updatedAt":  this.UpdatedAt.Unix(),
		"deletedAt":  this.DeletedAt,
	}
}

func (this *Collect) getVersions() []string {
	if this.Versions == nil {
		return []string{}
	}
	return this.Versions
}

func (this *Collect) Init() {
	this.AddFilters(transforms.FilterEmpty)
	this.SetProvider(DataProvider, this.data)
	this.SetProvider(SaverProvider, this.save)
	this.SetProvider(DefaultProvider, this.defaults)
	this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *Collect) defaults() {
	if this.Id == "" {
		this.Id = bson.NewObjectId()
	}
	if this.UpdatedAt.IsZero() {
		this.UpdatedAt = time.Now().Local()
	}
	if this.CreatedAt.IsZero() {
		this.CreatedAt = time.Now().Local()
	}
	if this.Status == 0 {
		this.Status = 1
	}
	if this.TargetType == "" {
		this.TargetType = CollectTargetTypePost
	}
}

func (this *Collect) save() error {
	var (
		model = CollectModelOf()
		data  = model.GetByUnique(this.data())
	)
	if data == nil {
		this.InitDefault()
		return model.Add(this)
	}
	return model.Update(bson.M{"_id": data.Id}, this.M(func(m beego.M) beego.M {
		delete(m, "id")
		delete(m, "createdAt")
		m["updatedAt"] = time.Now().Local()
		return m
	}))
}

func (this *Collect) Set(key string, v interface{}) *Collect {
	switch key {
	case "id":
		this.SetObjectId(&this.Id, v)
	case "userId":
		this.SetString(&this.UserId, v)
	case "targetId":
		this.SetString(&this.TargetId, v)
	case "targetType":
		this.SetString(&this.TargetType, v)
	case "status":
		this.SetNumInt(&this.Status, v)
	case "versions":
		this.SetStringArr(&this.Versions, v)
	case "createdAt":
		this.SetTime(&this.CreatedAt, v)
	case "updatedAt":
		this.SetTime(&this.UpdatedAt, v)
	case "deletedAt":
		this.SetNumIntN(&this.DeletedAt, v)
	}
	return this
}

func (this *Collect) setAttributes(data map[string]interface{}, safe ...bool) {
	for key, v := range data {
		if !safe[0] {
			if this.Excludes(key) {
				continue
			}
			if this.IsEmpty(v) {
				continue
			}
		}
		this.Set(key, v)
	}
}

func (this *CollectModel) TableName() string {
	return CollectTable
}

func (this *CollectModel) CreateIndex(force ...bool) {
	this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *CollectModel) getCreateIndexHandler() func(*mgo.Collection) {
	return func(doc *mgo.Collection) {
		this.logs(doc.EnsureIndex(mgo.Index{
			Key:    []string{"userId", "targetId", "targetType"},
			Unique: true,
			Sparse: false,
		}))
		this.logs(doc.EnsureIndexKey("status"))
	}
}

func (this *CollectModel) init() *CollectModel {
	this.Bind(this)
	this.Init()
	return this
}

func (this *CollectModel) GetByUnique(m beego.M) *Collect {
	if len(m) == 0 {
		return nil
	}
	var (
		err   error
		data  = NewCollect()
		query = beego.M{"userId": "", "targetId": "", "targetType": ""}
	)
	for key, _ := range query {
		v, ok := m[key]
		if !ok {
			return nil
		}
		str, ok := v.(string)
		if !ok || str == "" {
			return nil
		}
		query[key] = str
	}
	err = this.FindOne(query, data)
	if err != nil {
		return nil
	}
	return data
}

func (this *CollectModel) GetTravelNote(id string) *TravelNotes {
	if this == nil || id == "" {
		return nil
	}
	var (
		model = PostsModelOf()
		posts = NewTravelNotes()
		query = bson.M{
			"_id": id,
		}
	)
	if err := model.FindOne(query, posts); err == nil {
		return posts
	}
	return nil
}

func (this *CollectModel) GetTravelNotesByIds(ids []string) []*TravelNotes {
	if this == nil || ids == nil || len(ids) <= 0 {
		return nil
	}
	var list []bson.ObjectId
	for _, v := range ids {
		if v == "" {
			continue
		}
		list = append(list, bson.ObjectIdHex(v))
	}
	if len(list) <= 0 {
		return nil
	}
	var (
		model    = PostsModelOf()
		postsArr = make([]*TravelNotes, 0)
		query    = bson.M{
			"_id": bson.M{
				"$in": list,
			},
		}
		err = model.Gets(query, &postsArr)
	)
	if err == nil {
		return postsArr
	}
	log.Println("error", err)
	return nil
}
