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

// Tag 标签记录
type Tag struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`              // ID
		Name          string        `json:"name" bson:"name"`           // 标签名
		Alias         string        `json:"alias" bson:"alias"`         // 别名
		Group         string        `json:"group" bson:"group"`         // 分组名
		Comment       string        `json:"comment" bson:"comment"`     // 备注
		State         int           `json:"state" bson:"state"`         // 状态 0:初始状态,1:正常,2:删除
		Sort          int           `json:"sort" bson:"sort"`           // 排序 越大越靠前
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"` // 创建时间
		dataClassImpl `json:",omitempty" bson:",omitempty"`
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
				"sort":      this.Sort,
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

// Set setter
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
				this.SetNumInt(&this.State, v)
		case "sort":
				this.SetNumInt(&this.Sort, v)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v, true)
		}
		return this
}

func (this *Tag) save() error {
		var (
				newTag = NewTag()
				model  = TagsModelOf()
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
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
}

func (this *TagModel) TableName() string {
		return TagModelTableName
}

func TagsModelOf() *TagModel {
		var model = new(TagModel)
		model.Bind(model)
		model.Init()
		return model
}

func (this *TagModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *TagModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"name", "group"},
						Unique: true,
						Sparse: false,
				}))
				this.logs(doc.EnsureIndexKey("state"))
		}
}

func (this *TagModel) GetTagsByGroup(group ...string) []string {
		table := this.Document()
		defer this.destroy(table)
		var (
				strArr []string
				arr    = make([]*Tag, 2)
		)
		arr = arr[:0]
		if err := table.Find(bson.M{"group": group[0], "state": 1}).Sort("-sort").All(&arr); err == nil {
				for _, it := range arr {
						strArr = append(strArr, it.Name)
				}
		}
		return strArr
}

// 获取对应所有标签
func (this *TagModel) GetTags(group string) []Tag {
		var (
				arr = make([]Tag, 2)
		)
		arr = arr[:0]
		if err := this.NewQuery(bson.M{"group": group, "state": 1}).Sort("-sort").All(&arr); err == nil {
				return arr
		}
		return arr
}

// 批量添加更新
func (this *TagModel) Adds(items []map[string]interface{}) error {
		if len(items) == 0 {
				return ErrEmptyData
		}
		var result []interface{}
		for _, it := range items {
				var tag = this.GetByUnique(it)
				if tag != nil {
						_ = this.Update(bson.M{"_id": tag.Id}, it)
				} else {
						var tag = NewTag()
						tag.SetAttributes(it, false)
						tag.InitDefault()
						result = append(result, tag)
				}
		}
		if len(result) == 0 {
				return nil
		}
		if err := this.Inserts(result); err != nil {
				return err
		}
		return nil
}

// GetByUnique 通过唯一索引查询
func (this *TagModel) GetByUnique(data map[string]interface{}) *Tag {
		var (
				info        = NewTag()
				name, group = data["name"], data["group"]
		)
		if name == nil || group == nil {
				return nil
		}
		err := this.FindOne(bson.M{"name": name, "group": group}, info)
		if err == nil {
				return info
		}
		return nil
}
