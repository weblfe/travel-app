package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/libs"
		"time"
)

type MessageTemplateModel struct {
		BaseModel
}

type MessageTemplate struct {
		Id         bson.ObjectId `json:"id" bson:"_id"`                                    // 目标Id
		Type       string        `json:"type" bson:"type"`                                 // 模版类型名
		Title      string        `json:"title" bson:"title"`                               // 模版标题
		Name       string        `json:"name" bson:"name"`                                 // 模版名称
		Template   bson.M        `json:"template" bson:"template"`                         // 模版信息
		TemplateId string        `json:"templateId,omitempty" bson:"templateId,omitempty"` // 第三方模版ID
		Platform   string        `json:"platform" bson:"platform"`                         // 平台
		Comment    string        `json:"comment,omitempty" bson:"comment,omitempty"`       // 备注
		State      int           `json:"state" bson:"state"`                               // 状态 0:不可用,1:可用
		UpdatedAt  time.Time     `json:"updatedAt" bson:"updatedAt"`                       // 更新时间
		CreatedAt  time.Time     `json:"createdAt" bson:"createdAt"`                       // 创建时间
}

// 短信消息模版
type SmsTemplate struct {
		Content   string `json:"content"`
		Variables []*struct {
				Key   string `json:"key"`
				Value string `json:"value"`
				Desc  string `json:"desc"`
		} `json:"variables"`
}

const (
		MessageTemplateModelTableName = "message_template"
)

func NewMessageTemplate() *MessageTemplate {
		var template = new(MessageTemplate)
		template.State = 1
		return template
}

func MessageTemplateModelOf() *MessageTemplateModel {
		var model = new(MessageTemplateModel)
		model.Bind(model)
		model.Init()
		return model
}

func NewSmsTemplate() *SmsTemplate {
		return new(SmsTemplate)
}

func (this *SmsTemplate) M() beego.M {
		return beego.M{
				"content":   this.Content,
				"variables": this.variables(),
		}
}

func (this *SmsTemplate) Load(data map[string]interface{}) *SmsTemplate {
		if str, err := libs.Json().Marshal(data); err == nil {
				_ = libs.Json().Unmarshal(str, this)
		}
		return this
}

func (this *SmsTemplate) variables() []beego.M {
		var arr []beego.M
		for _, it := range this.Variables {
				arr = append(arr, beego.M{
						"key":   it.Key,
						"value": it.Value,
						"desc":  it.Desc,
				})
		}
		return arr
}

func (this *MessageTemplate) Load(data map[string]interface{}) *MessageTemplate {
		for key, v := range data {
				this.Set(key, v)
		}
		return this
}

func (this *MessageTemplate) Set(key string, v interface{}) *MessageTemplate {
		switch key {
		case "comment":
				this.Comment = v.(string)
		case "title":
				this.Title = v.(string)
		case "template":
				if m, ok := v.(bson.M); ok {
						this.Template = m
				}
				if m, ok := v.(beego.M); ok {
						this.Template = bson.M(m)
				}
				if m, ok := v.(map[string]interface{}); ok {
						this.Template = m
				}
				if m, ok := v.(*SmsTemplate); ok {
						this.Template = bson.M(m.M())
				}
		case "type":
				this.Type = v.(string)
		case "state":
				this.State = v.(int)
		case "name":
				this.Name = v.(string)
		case "templateId":
				this.TemplateId = v.(string)
		case "platform":
				this.Platform = v.(string)
		case "updatedAt":
				this.UpdatedAt = v.(time.Time)
		case "createdAt":
				this.CreatedAt = v.(time.Time)
		}
		return this
}

func (this *MessageTemplate) Defaults() *MessageTemplate {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now().Local()
		}
		if this.State == 0 {
				this.State = 1
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		return this
}

func (this *MessageTemplateModel) TableName() string {
		return MessageTemplateModelTableName
}

// 创建索引
func (this *MessageTemplateModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *MessageTemplateModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"type", "name"},
						Unique: true,
						Sparse: true,
				}))
				this.logs(doc.EnsureIndexKey("state"))
				this.logs(doc.EnsureIndexKey("type"))
				this.logs(doc.EnsureIndexKey("template"))
				this.logs(doc.EnsureIndexKey("platform"))
		}
}

// 批量添加
func (this *MessageTemplateModel) Adds(data []map[string]interface{}) error {
		for _, it := range data {
				tmp := NewSmsTemplate()
				v, ok := it["template"]
				if !ok {
						continue
				}
				if m, ok := v.(map[string]interface{}); ok {
						tmp.Load(m)
						it["template"] = tmp
				}
				t := NewMessageTemplate()
				t.Load(it)
				query := bson.M{
						"name":       t.Name,
						"type":       t.Type,
						"title":      t.Title,
						"templateId": t.TemplateId,
						"platform":   t.Platform,
				}

				template := this.GetByUnique(query)
				if template != nil {
						_ = this.Update(bson.M{"_id": template.Id}, it)
						continue
				}
				t = t.Defaults()
				if err := this.Add(t); err != nil {
						return err
				}
		}
		return nil
}

// 通过唯一查询条件获取
func (this *MessageTemplateModel) GetByUnique(data map[string]interface{}) *MessageTemplate {
		var (
				name, typ, title, templateId, platform = data["name"], data["type"], data["title"], data["templateId"], data["platform"]
		)
		if name == "" || name == nil || typ == nil || templateId == nil || templateId == "" || platform == "" {
				return nil
		}
		query := bson.M{
				"name":       name,
				"type":       typ,
				"title":      title,
				"templateId": templateId,
				"platform":   platform,
		}
		var info = NewMessageTemplate()
		if err := this.FindOne(query, info); err == nil {
				return info
		}
		return nil
}
