package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"time"
)

// 审核日志
type AuditLog struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`              // ID
		UserId        string        `json:"userId" json:"userId"`       // 审核用户ID
		Platform      int           `json:"platform" bson:"platform"`   // 审核端 [1:ios,2,android,3,pc]
		PostId        string        `json:"postId" bson:"postId"`       // 文件ID
		AuditType     string        `json:"auditType" bson:"auditType"` // 类型
		Comment       string        `json:"comment" bson:"comment"`     // 备注
		UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"` // 更新时间时间
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"` // 创建时间
		dataClassImpl `bson:",omitempty"  json:",omitempty"`
}

type AuditLogModel struct {
		BaseModel
}

const (
		AuditLogTable = "audit_logs"
)

func NewAuditLog() *AuditLog {
		return new(AuditLog)
}

func AuditLogModelOf() *AuditLogModel {
		var model = new(AuditLogModel)
		model.Bind(model)
		model.init()
		return model
}

func (this *AuditLog) Init() {
		//	this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *AuditLog) data() beego.M {
		return beego.M{
				"id":        this.Id.Hex(),
				"platform":  this.Platform,
				"userId":    this.UserId,
				"postId":    this.PostId,
				"auditType": this.AuditType,
				"comment":   this.Comment,
				"updatedAt": this.UpdatedAt.Unix(),
				"createdAt": this.CreatedAt.Unix(),
		}
}

func (this *AuditLog) save() error {
		var (
				info  = NewAuditLog()
				model = AuditLogModelOf()
		)
		err := model.FindOne(bson.M{"platform": this.Platform, "postId": this.PostId, "auditType": this.AuditType}, info)
		if err == nil {
				info.setAttributes(this.M())
				return model.Update(bson.M{"_id": info.Id}, info)
		}
		return model.Add(this)
}

func (this *AuditLogModel) GetByUnique(m beego.M) *AuditLog {
		var (
				info  = NewAuditLog()
				query = beego.M{}
		)
		if len(query) == 0 {
				return nil
		}
		var keys = []string{"platform", "postId", "auditType"}
		for _, key := range keys {
				v, ok := m[key]
				if !ok {
						return nil
				}
				query[key] = v
		}
		err := this.FindOne(query, info)
		if err == nil {
				return info
		}
		return nil
}

func (this *AuditLog) setAttributes(data map[string]interface{}, safe ...bool) {
		if len(safe) == 0 {
				safe = append(safe, false)
		}
		for k, v := range data {
				if safe[0] {
						// 排除键
						if this.Excludes(k) {
								continue
						}
						if this.IsEmpty(v) {
								continue
						}
				}
				this.Set(k, v)
		}
}

func (this *AuditLog) Set(key string, v interface{}) *AuditLog {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, v)
		case "platform":
				this.SetNumInt(&this.Platform, v)
		case "comment":
				this.SetString(&this.Comment, v)
		case "postId":
				this.SetString(&this.PostId, v)
		case "userId":
				this.SetString(&this.UserId, v)
		case "auditType":
				this.SetString(&this.AuditType, v)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
		case "updatedAt":
				this.SetTime(&this.UpdatedAt, v)
		}
		return this
}

func (this *AuditLog) setDefaults() {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now().Local()
		}
}

func (this *AuditLogModel) init() {
		this.Init()
}

func (this *AuditLogModel) TableName() string {
		return AuditLogTable
}

func (this *AuditLogModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *AuditLogModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				// unique mobile
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"userId", "postId", "platform", "auditType", "createdAt"},
						Unique: true,
						Sparse: false,
				}))
		}
}
