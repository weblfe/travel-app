package models

import (
		"encoding/json"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"strconv"
		"time"
)

type MessageModel struct {
		BaseModel
}

type MessageLog struct {
		Id             bson.ObjectId `json:"id" bson:"id"`                                             // 消息ID
		Title          string        `json:"title" bson:"title"`                                       // 消息标题
		Type           string        `json:"type" bson:"type"`                                         // 消息类型 [register.sms,register.email,reset.sms,rebind.email,bind.sms]
		SenderProvider string        `json:"provider" bson:"provider"`                                 // 发送服务名
		Extras         bson.M        `json:"extras,omitempty" bson:"extras,omitempty"`                 // 扩展信息
		Content        string        `json:"content" bson:"content"`                                   // 消息内容
		SenderUserId   string        `json:"sender_user_id,omitempty" bson:"sender_user_id,omitempty"` // 发送人
		TargetUserId   string        `json:"target_user_id,omitempty" bson:"target_user_id,omitempty"` // 接受人
		Mobile         string        `json:"mobile,omitempty" bson:"mobile,omitempty"`                 // 手机消息手机号
		Email          string        `json:"email,omitempty" bson:"email,omitempty"`                   // 邮箱信息邮箱号
		State          int           `json:"state" bson:"state"`                                       // 消息状态 [-3:拒绝接收,-2:发送失败,-1:待处理,0:未知,1:已发送,2:已阅读]
		Result         string        `json:"result,omitempty" bson:"result,omitempty"`                 // 第三方消息结果
		SentTime       int64         `json:"sent_time" bson:"sent_time"`                               // 发送时间
		ExpireTime     int64         `json:"expire_time,omitempty" bson:"expire_time,omitempty"`       // 消息过期时间
		ReadTime       int64         `json:"read_time,omitempty" bson:"read_time,omitempty"`           // 消息阅读时间
		CreatedAt      time.Time     `json:"created_at" bson:"created_at"`                             // 记录创建时间
}

const (
		MessageLogTable = "message_log"
)

func NewMessageLog() *MessageLog {
		var message = new(MessageLog)
		return message
}

func MessageModelOf() *MessageModel {
		var model = new(MessageModel)
		model._Self = model
		model.Init()
		return model
}

func (this *MessageLog) Load(data map[string]interface{}) *MessageLog {
		for k, v := range data {
				this.Set(k, v)
		}
		return this
}

func (this *MessageLog) Defaults() *MessageLog {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now()
		}
		if this.State == 0 {
				this.State = -1
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		return this
}

func (this *MessageLog) Set(key string, v interface{}) *MessageLog {
		switch key {
		case "title":
				this.Title = v.(string)
		case "type":
				this.Type = v.(string)
		case "provider":
				this.SenderProvider = v.(string)
		case "extras":
				if m, ok := v.(bson.M); ok {
						this.Extras = m
				}
				if m, ok := v.(beego.M); ok {
						this.Extras = bson.M(m)
				}
				if m, ok := v.(map[string]interface{}); ok {
						this.Extras = m
				}
		case "content":
				this.Content = v.(string)
		case "sender_user_id":
				this.SenderUserId = v.(string)
		case "target_user_id":
				this.TargetUserId = v.(string)
		case "mobile":
				this.Mobile = v.(string)
		case "email":
				this.Email = v.(string)
		case "state":
				if s, ok := v.(string); ok {
						if n, err := strconv.Atoi(s); err == nil {
								this.State = n
						}
						return this
				}
				this.State = v.(int)
		case "result":
				if str, ok := v.(fmt.Stringer); ok {
						this.Result = str.String()
						return this
				}
				if str, ok := v.(string); ok {
						this.Result = str
						return this
				}
				b, _ := json.Marshal(v)
				this.Result = string(b)
		case "sent_time":
				this.SentTime = v.(int64)
		case "read_time":
				this.ReadTime = v.(int64)
		case "expire_time":
				this.ExpireTime = v.(int64)
		case "created_at":
				this.CreatedAt = v.(time.Time)
		}
		return this
}

func (this *MessageModel) CreateIndex() {
		_ = this.Collection().EnsureIndexKey("state")
		_ = this.Collection().EnsureIndexKey("mobile")
		_ = this.Collection().EnsureIndexKey("type")
		_ = this.Collection().EnsureIndexKey("provider")
		_ = this.Collection().EnsureIndexKey("email")
		_ = this.Collection().EnsureIndexKey("sent_time")
		_ = this.Collection().EnsureIndexKey("expire_time")
		_ = this.Collection().EnsureIndexKey("created_at")
		_ = this.Collection().EnsureIndexKey("sender_user_id", "target_user_id")
		_ = this.Collection().EnsureIndexKey("extras")
}

func (this *MessageModel) TableName() string {
		return MessageLogTable
}
