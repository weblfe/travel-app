package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"time"
)

type AppModel struct {
		BaseModel
}

const (
		DefaultAppName  = "游迹"
		DefaultVersion  = "1.0.0"
		AppVersionTable = "app_versions_info"
)

func AppModelOf() *AppModel {
		var model = new(AppModel)
		model._Self = model
		model.Init()
		return model
}

// app Built 信息模型
type AppInfo struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`                  // ID
		Driver        string        `json:"driver" bson:"driver"`           // 设备类型 ios,android,win
		Download      string        `json:"download" bson:"download"`       // 下载链接
		Version       string        `json:"version" bson:"version"`         // 版本号
		Content       string        `json:"content" bson:"content"`         // 更新内容
		ForcedUpdate  bool          `json:"state" bson:"forcedUpdate"`      // 是否强制更新
		PublishTime   int64         `json:"publishTime" bson:"publishTime"` // 版本发布时间
		Remark        string        `json:"remark" bson:"remark"`           // 备注
		AppName       string        `json:"appName" bson:"appName"`         // 应用名
		AppBuild      int           `json:"appBuild" bson:"appBuild"`       // 构建版本数字版本
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`     // 创建时间
		dataClassImpl `bson:",omitempty"  json:",omitempty"`
}

func NewAppInfo() *AppInfo {
		var info = new(AppInfo)
		info.Init()
		return info
}

func (this *AppInfo) Init() {
		//	this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *AppInfo) data() beego.M {
		return beego.M{
				"id":           this.Id.Hex(),
				"driver":       this.Driver,
				"download":     this.Download,
				"version":      this.Version,
				"content":      this.Content,
				"forcedUpdate": this.ForcedUpdate,
				"remark":       this.Remark,
				"appName":      this.GetAppName(),
				"appBuild":     this.GetAppBuild(),
				"publishTime":  this.PublishTime,
				"createdAt":    this.CreatedAt.Unix(),
		}
}

func (this *AppInfo) GetAppName() string {
		if this.AppName != "" {
				return this.AppName
		}
		var name = ConfigModelOf().GetString("appName")
		if name == "" {
				return DefaultAppName
		}
		return name
}

func (this *AppInfo) GetAppBuild() int {
		if this.AppBuild != 0 {
				return this.AppBuild
		}
		return AppModelOf().GetAppBuild(this.Driver) + 1
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

func (this *AppInfo) Set(key string, v interface{}) *AppInfo {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, v)
		case "driver":
				this.SetString(&this.Driver, v)
		case "download":
				this.SetString(&this.Download, v)
		case "content":
				this.SetString(&this.Content, v)
		case "version":
				this.SetString(&this.Version, v)
		case "state":
				this.SetBool(&this.ForcedUpdate, v)
		case "forcedUpdate":
				this.SetBool(&this.ForcedUpdate, v)
		case "remark":
				this.SetString(&this.Remark, v)
		case "appName":
				this.SetString(&this.AppName, v)
		case "appBuild":
				this.SetNumInt(&this.AppBuild, v)
		case "publishTime":
				if str, ok := v.(string); ok {
						t, err := time.Parse("2020/10/07 10:00:00", str)
						if err != nil {
								return this
						}
						this.PublishTime = t.Unix()
						return this
				}
				this.PublishTime = v.(int64)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
		}
		return this
}

func (this *AppInfo) setDefaults() {

		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.PublishTime == 0 {
				this.PublishTime = time.Now().Unix()
		}
		if this.AppBuild == 0 {
				this.AppBuild = this.GetAppBuild()
		}
		if this.AppName == "" {
				this.AppName = this.GetAppName()
		}
		if this.Download == "" {
				this.Download = this.GetDownloadUrl()
		}
}

func (this *AppInfo) GetDownloadUrl() string {
		if this.Download != "" {
				return this.Download
		}
		return AppModelOf().GetAppDownloadUrl(this.Driver)
}

// 表名
func (this *AppModel) TableName() string {
		return AppVersionTable
}

// 创建索引
func (this *AppModel) CreateIndex() {
		// unique mobile
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"driver", "version"},
				Unique: true,
				Sparse: false,
		})
		_ = this.Collection().EnsureIndexKey("forcedUpdate")
}

// 批量添加更新
func (this *AppModel) Adds(items []map[string]interface{}) error {
		if len(items) == 0 {
				return ErrEmptyData
		}
		var result []interface{}
		for _, it := range items {
				info := this.GetByUnique(it)
				if info != nil {
						info.Init()
						info.setAttributes(it)
						info.InitDefault()
						_ = this.Update(bson.M{"_id": info.Id}, info.M(info.GetFormatterTime("createdAt")))
				} else {
						info := NewAppInfo()
						info.SetAttributes(it, false)
						info.InitDefault()
						result = append(result, info)
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

// 获取最新版本号
func (this *AppModel) GetAppBuild(driver string) int {
		var (
				err  error
				info = NewAppInfo()
		)
		err = this.NewQuery(bson.M{"driver": driver}).Sort("-appBuild").One(info)
		if err == nil {
				return info.AppBuild
		}
		return 0
}

// 获取最新版本号
func (this *AppModel) GetAppVersion(driver string) string {
		var (
				err  error
				info = NewAppInfo()
		)
		err = this.NewQuery(bson.M{"driver": driver}).Sort("-appBuild").One(info)
		if err == nil {
				return info.Version
		}
		return DefaultVersion
}

// 获取最新版本号
func (this *AppModel) GetAppDownloadUrl(driver string) string {
		var (
				err  error
				info = NewAppInfo()
		)
		err = this.NewQuery(bson.M{"driver": driver}).Sort("-appBuild").One(info)
		if err == nil {
				return info.Download
		}
		return ""
}

// 通过唯一索引查询
func (this *AppModel) GetByUnique(data map[string]interface{}) *AppInfo {
		var (
				info            = NewAppInfo()
				driver, version = data["driver"], data["version"]
		)
		if driver == nil || version == nil {
				return nil
		}
		err := this.FindOne(bson.M{"driver": driver, "version": version}, info)
		if err == nil {
				return info
		}
		return nil
}
