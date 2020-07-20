package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/transforms"
		"time"
)

type TagModel struct {
		BaseModel
}

// 标签记录
type Tag struct {
		Id        bson.ObjectId `json:"id" bson:"_id"`              // ID
		Name      string        `json:"name" bson:"name"`           // 标签名
		Alias     string        `json:"alias" bson:"alias"`         // 别名
		Group     string        `json:"group" bson:"group"`         // 分组名
		Comment   string        `json:"comment" bson:"comment"`     // 备注
		State     int           `json:"state" bson:"state"`         // 状态 0:初始状态,1:正常,2:删除
		CreatedAt time.Time     `json:"createdAt" bson:"createdAt"` // 创建时间
		dataClassImpl
}

const (
		TagModelTableName = "tags"
)

func NewTag() *Tag {
		var tag = new(Tag)
		tag.Init()
		return tag
}

func (this *Tag) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *Tag) data() beego.M {
		return beego.M{
				"id":        this.Id.Hex(),
				"name":      this.Name,
				"comment":   this.Comment,
				"alias":     this.Alias,
				"group":     this.Group,
				"state":     this.State,
				"createdAt": this.CreatedAt.Unix(),
		}
}

func (this *Tag) setAttributes(data map[string]interface{}, safe ...bool) {
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

func (this *Tag) Set(key string, v interface{}) *Tag {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, v)
		case "name":
				this.SetString(&this.Name, v)
		case "group":
				this.SetString(&this.Group, v)
		case "comment":
				this.SetString(&this.Comment, v)
		case "state":
				this.State = v.(int)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v, true)
		}
		return this
}

func (this *Tag) save() error {
		var (
				newTag = NewTag()
				model  = TagModelOf()
		)
		err := model.FindOne(beego.M{"name": this.Name, "group": this.Group}, newTag)
		if err == nil {
				newTag.setAttributes(this.M())
				return model.Update(beego.M{"_id": newTag.Id}, newTag)
		}
		this.InitDefault()
		return model.Update(beego.M{"_id": this.Id}, this)
}

func (this *Tag) setDefaults() {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.State == 0 {
				this.State = 1
		}
		if this.Group == "" {
				this.Group = "post"
		}
}

func (this *TagModel) TableName() string {
		return TagModelTableName
}

func TagModelOf() *TagModel {
		var model = new(TagModel)
		model._Self = model
		model.Init()
		return model
}

func (this *TagModel) CreateIndex() {
		// unique mobile
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"name", "group"},
				Unique: true,
				Sparse: false,
		})
		_ = this.Collection().EnsureIndexKey("state")
}

func (this *TagModel) GetTagsByGroup(group ...string) []string {
		table := this.Collection()
		defer this.destroy()
		var (
				strArr []string
				arr    = make([]*Tag, 2)
		)
		arr = arr[:0]
		if err := table.Find(bson.M{"group": group[0], "state": 1}).All(&arr); err == nil {
				for _, it := range arr {
						strArr = append(strArr, it.Name)
				}
		}
		return strArr
}