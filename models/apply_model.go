package models

import (
		"errors"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"time"
)

// 工单 （申请｜举报） 记录
type ApplyInfo struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`              // 工单ID
		Title         string        `json:"title" bson:"title"`         // 标题
		UserId        string        `json:"userId" bson:"userId"`       // 申请｜举报人
		Type          string        `json:"type" bson:"type"`           // 类型
		Target        string        `json:"target" bson:"target"`       // 举报｜申请目标 | mobile
		Images        []string      `json:"images" bson:"images"`       // 图片集合
		Status        int           `json:"status" bson:"status"`       // 状态
		Content       string        `json:"content" bson:"content"`     // 内容
		Extras        beego.M       `json:"extras" bson:"extras"`       // 扩展
		Date          int64         `json:"date" bson:"date"`           // 日期
		UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"` // 更新时间
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"` // 创建时间
		dataClassImpl `bson:",omitempty"  json:",omitempty"`
}

type ApplyInfoModel struct {
		BaseModel
}

const (
		ApplyInfoTable     = "apply_infos" // 表名
		ApplyTypeReport    = "report"      // 举报
		ApplyTypeSuggest   = "suggest"     // 建议｜反馈
		ApplyTypeAllege    = "allege"      // 申诉
		ApplyStatusWait    = 0             // 待处理
		ApplyStatusPassed  = 1             // 处理完毕
		ApplyStatusDealing = 2             // 处理中
		ApplyStatusReject  = 3             // 拒绝处理
)

var (
		ApplyTypes = []string{
				ApplyTypeReport,
				ApplyTypeSuggest,
				ApplyTypeAllege,
		}
		ApplyStatusDesc = map[int]string{
				ApplyStatusWait:    "待处理",
				ApplyStatusPassed:  "处理完毕",
				ApplyStatusDealing: "处理中",
				ApplyStatusReject:  "拒绝处理",
		}
)

func NewApplyInfo() *ApplyInfo {
		var info = new(ApplyInfo)
		info.init()
		return info
}

func (this *ApplyInfo) init() *ApplyInfo {
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.defaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
		return this
}

func (this *ApplyInfo) data() beego.M {
		return beego.M{
				"id":         this.Id.Hex(),
				"userId":     this.UserId,
				"title":      this.Title,
				"type":       this.Type,
				"target":     this.Target,
				"content":    this.Content,
				"extra":      this.Extras,
				"status":     this.Status,
				"date":       this.Date,
				"images":     this.getImages(),
				"statusDesc": this.getStatusText(),
				"updatedAt":  this.UpdatedAt.Unix(),
				"createdAt":  this.CreatedAt.Unix(),
		}
}

func (this *ApplyInfo) getStatusText() string {
		var desc = ApplyStatusDesc[this.Status]
		return desc
}

// 获取图片
func (this *ApplyInfo) getImages() []*Image {
		var images = make([]*Image, 2)
		if this.Images == nil || len(this.Images) <= 0 {
				return images[:0]
		}
		var attach = AttachmentModelOf()
		StrArray(this.Images).Foreach(func(i int, value string) bool {
				if value == "" || i < 0 {
						return true
				}
				img := attach.GetImageById(value)
				if img != nil {
						images = append(images, img)
				}
				return true
		})
		return images
}

func (this *ApplyInfo) CheckType() bool {
		return StrArray(ApplyTypes).Included(this.Type)
}

// 验证图片是否可用
func (this *ApplyInfo) AppendImages(images ...string) *ApplyInfo {
		var attach = AttachmentModelOf()
		StrArray(images).Foreach(func(i int, value string) bool {
				if value == "" || i < 0 {
						return true
				}
				img := attach.GetImageById(value)
				if img != nil {
						this.Images = append(this.Images, value)
				}
				return true
		})
		return this
}

func (this *ApplyInfo) Set(key string, value interface{}) *ApplyInfo {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, value)
		case "userId":
				this.SetString(&this.UserId, value)
		case "title":
				this.SetString(&this.Title, value)
		case "type":
				this.SetString(&this.Type, value)
		case "target":
				this.SetString(&this.Target, value)
		case "images":
				this.SetStringArr(&this.Images, value)
		case "content":
				this.SetString(&this.Content, value)
		case "extra":
				this.SetMapper(&this.Extras, value)
		case "status":
				this.SetNumInt(&this.Status, value)
		case "updatedAt":
				this.SetTime(&this.UpdatedAt, value)
		case "createdAt":
				this.SetTime(&this.CreatedAt, value)
		case "date":
				this.SetNumIntN(&this.Date, value)
		}
		return this
}

func (this *ApplyInfo) save() error {
		var (
				id    = this.Id.Hex()
				tmp   = NewApplyInfo()
				model = ApplyInfoModelOf()
				err   = model.GetById(id, tmp)
		)
		if err == nil {
				return model.UpdateById(id, this.M(func(m beego.M) beego.M {
						m = this.removeUpdateExcludes(m)
						m["updatedAt"] = time.Now().Local()
						return m
				}))
		}
		if this.Type == "" || this.Content == "" {
				return errors.New("params miss")
		}
		this.InitDefault()
		return model.Add(this)
}

func (this *ApplyInfo) removeUpdateExcludes(m beego.M) beego.M {
		var arr = []string{"createdAt", "updatedAt"}
		for _, key := range arr {
				delete(m, key)
		}
		return m
}

func (this *ApplyInfo) defaults() {
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now().Local()
		}
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.Extras == nil {
				this.Extras = beego.M{}
		}
		if this.Date == 0 {
				this.Date = GetDate()
		}
}

func (this *ApplyInfo) setAttributes(data map[string]interface{}, safe ...bool) {
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

func ApplyInfoModelOf() *ApplyInfoModel {
		var model = new(ApplyInfoModel)
		model.Bind(model)
		model.Init()
		return model
}

func (this *ApplyInfoModel) TableName() string {
		return ApplyInfoTable
}

func (this *ApplyInfoModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *ApplyInfoModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				// unique mobile
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"userId", "target", "type", "title"},
						Unique: true,
						Sparse: false,
				}))
				this.logs(doc.EnsureIndexKey("date"))
				this.logs(doc.EnsureIndexKey("status"))
		}
}

func (this *ApplyInfoModel) Count(query bson.M) int {
		var count, err = this.NewQuery(query).Count()
		if err == nil {
				return count
		}
		return 0
}

func (this *ApplyInfoModel) GetByUnique(m beego.M) *ApplyInfo {
		var (
				err   error
				data  = NewApplyInfo()
				query = bson.M{"userId": "", "target": "", "type": "", "title": "","content":""}
		)
		for key := range query {
				v, ok := m[key]
				if !ok {
						return nil
				}
				if key != "date" {
						str, ok := v.(string)
						if ok && str != "" {
								query[key] = str
								continue
						}
				} else {
						str, ok := v.(string)
						if ok && str != "" {
								t, err := time.Parse(time.RFC3339, str)
								if err != nil {
										return nil
								}
								query[key] = t.Unix()
								continue
						}
						t, ok := v.(int64)
						if ok && str != "" {
								query[key] = t
								continue
						}
				}
				return nil
		}
		err = this.NewQuery(query).One(data)
		if err == nil {
				return data
		}
		return nil
}
