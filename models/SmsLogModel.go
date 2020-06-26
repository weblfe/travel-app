package models

import (
		"github.com/globalsign/mgo/bson"
		"time"
)

type SmsLogModel struct {
		BaseModel
}

// 短信日志记录
type SmsLog struct {
		Id              bson.ObjectId `json:"id" bson:"id"`                               // 消息ID
		Provider        string        `json:"provider" bson:"provider"`                   // 服务商
		Mobile          string        `json:"mobile" bson:"mobile"`                       // 手机号
		Content         string        `json:"content" bson:"content"`                     // 消息内容
		Result          string        `json:"result" bson:"result"`                       // 请求结果
		State           int           `json:"state" bson:"state"`                         // 消息状态
		Type            string        `json:"type" bson:"type"`                           // 消息类型
		Extras          string        `json:"extras" bson:"extras"`                       // 扩展信息
		Error           string        `json:"error,omitempty" bson:"error,omitempty"`     // 异常
		ExpireTimeStamp int64         `json:"expire_time_stamp" bson:"expire_time_stamp"` // 过期时间戳
		CreatedAt       time.Time     `json:"created_at" bson:"created_at"`               // 创建时间
}

const (
		SmsLogTableName     = "sms_log"
		DefaultProviderName = "aliyun-dysms"
)

func SmsLogModelOf() *SmsLogModel {
		var model = new(SmsLogModel)
		model._Self = model
		model.Init()
		return model
}

func (this *SmsLogModel) CreateIndex() {
		_ = this.Collection().EnsureIndexKey("state")
		_ = this.Collection().EnsureIndexKey("type")
		_ = this.Collection().EnsureIndexKey("mobile")
		_ = this.Collection().EnsureIndexKey("provider")
		_ = this.Collection().EnsureIndexKey("created_at")
}

func (this *SmsLogModel) TableName() string {
		return SmsLogTableName
}

func (this *SmsLog) Load(data map[string]interface{}) *SmsLog {
		for key, it := range data {
				this.Set(key, it)
		}
		return this
}

func (this *SmsLog) Set(key string, v interface{}) *SmsLog {
		switch key {
		case "provider":
				this.Provider = v.(string)
		case "mobile":
				this.Mobile = v.(string)
		case "content":
				this.Content = v.(string)
		case "result":
				this.Result = v.(string)
		case "state":
				this.State = v.(int)
		case "type":
				this.Type = v.(string)
		case "extras":
				this.Extras = v.(string)
		case "expire_time_stamp":
				this.ExpireTimeStamp = v.(int64)
		case "created_at":
				this.CreatedAt = v.(time.Time)
		case "id":
				if id, ok := v.(bson.ObjectId); ok {
						this.Id = id
				}
		}
		return this
}

func (this *SmsLog) Defaults() *SmsLog {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now()
		}
		if this.Provider == "" {
				this.Provider = DefaultProviderName
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		return this
}
