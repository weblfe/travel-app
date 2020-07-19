package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/transforms"
		"time"
)

type AppModel struct {
		BaseModel
}

const (
		AppVersionTable = "app_version_info"
)

func AppModelOf() *AppModel {
		var model = new(AppModel)
		model._Self = model
		model.Init()
		return model
}

// app 信息模型
type AppInfo struct {
		Id        bson.ObjectId `json:"id" bson:"_id"`              // ID
		Driver    string        `json:"driver" bson:"driver"`       // 设备类型 ios,android,win
		Download  string        `json:"download" bson:"download"`   // 下载链接
		Version   string        `json:"version" bson:"version"`     // 版本号
		Content   string        `json:"content" bson:"content"`     // 更新内容
		State     int           `json:"state" bson:"state"`         // 状态
		CreatedAt time.Time     `json:"createdAt" bson:"createdAt"` // 创建时间
		dataClassImpl
}

func NewAppInfo() *AppInfo {
		var info = new(AppInfo)
		info.Init()
		return info
}

func (this *AppInfo) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *AppInfo) data() beego.M {
		return beego.M{
				"id":        this.Id.Hex(),
				"driver":    this.Driver,
				"download":  this.Download,
				"version":   this.Version,
				"content":   this.Content,
				"state":     this.State,
				"createdAt": this.CreatedAt.Unix(),
		}
}

func (this *AppInfo) save() error {
		var (
				info  = NewAppInfo()
				model = AppModelOf()
		)
		err := model.FindOne(bson.M{"version": this.Version, "driver": this.Driver}, info)
		if err == nil {
				info.setAttributes(this.M())
				return model.Update(bson.M{"_id": info.Id}, info)
		}
		return model.Add(this)
}

func (this *AppInfo) setAttributes(data map[string]interface{}, safe ...bool) {
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

func (this *AppInfo) Set(key string, v interface{}) *AppInfo {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, v)
		case "driver":
				this.SetString(&this.Driver, v)
		case "download":
				this.SetString(&this.Download, v)
		case "version":
				this.SetString(&this.Version, v)
		case "state":
				this.State = v.(int)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
		}
		return this
}

func (this *AppInfo) setDefaults() {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.State == 0 {
				this.State = 1
		}
}

func (this *AppModel) TableName() string {
		return AppVersionTable
}

func (this *AppModel) CreateIndex() {
		// unique mobile
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"driver", "version"},
				Unique: true,
				Sparse: false,
		})
		_ = this.Collection().EnsureIndexKey("state")
}
