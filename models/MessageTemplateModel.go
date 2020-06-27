package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"time"
)

type MessageTemplateModel struct {
		BaseModel
}

type MessageTemplate struct {
		Id        bson.ObjectId `json:"id" bson:"id"`                 // 目标Id
		Type      string        `json:"type" bson:"type"`             // 模版类型名
		Name      string        `json:"name" bson:"name"`             // 模版名称
		Template  bson.M        `json:"template" bson:"template"`     // 模版信息
		Platform  string        `json:"platform" bson:"platform"`     // 平台
		Comment   string        `json:"comment" bson:"comment"`       // 备注
		State     int           `json:"state" bson:"state"`           // 状态 0:不可用,1:可用
		UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"` // 更新时间
		CreatedAt time.Time     `json:"created_at" bson:"created_at"` // 创建时间
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
		model._Self = model
		model.Init()
		return model
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
		case "type":
				this.Type = v.(string)

		case "state":
				this.State = v.(int)
		case "updated_at":
				this.UpdatedAt = v.(time.Time)
		case "created_at":
				this.CreatedAt = v.(time.Time)
		}
		return this
}

func (this *MessageTemplate) Defaults() *MessageTemplate {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now()
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

func (this *MessageTemplateModel) CreateIndex() {
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"type", "name"},
				Unique: true,
				Sparse: true,
		})
		_ = this.Collection().EnsureIndexKey("state")
		_ = this.Collection().EnsureIndexKey("type")
		_ = this.Collection().EnsureIndexKey("template")
		_ = this.Collection().EnsureIndexKey("platform")
}
