package models

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/transforms"
		"strings"
		"time"
)

type ConfigModel struct {
		BaseModel
}

type Config struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`              // ID
		Key           string        `json:"key" bson:"key"`             // 键名
		Value         interface{}   `json:"value" json:"value"`         // 值
		Root          string        `json:"root" bson:"root"`           // 分组
		State         int           `json:"state" bson:"state"`         // 状态
		Title         string        `json:"title" bson:"title"`         // 备注
		UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"` // 更新时间
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"` // 创建时间
		dataClassImpl `bson:",omitempty"  json:",omitempty"`
}

const (
		ConfigTable       = "configs"
		DefaultConfigRoot = "app"
		ConfigStateUnInit = 0
		ConfigStateOk     = 1
		ConfigStateDel    = 2
)

// 配置模型
func ConfigModelOf() *ConfigModel {
		var model = new(ConfigModel)
		model.Bind(model)
		model.Init()
		return model
}

func NewConfig() *Config {
		var config = new(Config)
		config.Init()
		return config
}

func (this *Config) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *Config) data() beego.M {
		return beego.M{
				"id":        this.Id.Hex(),
				"key":       this.Key,
				"value":     this.Value,
				"root":      this.Root,
				"state":     this.State,
				"title":     this.Title,
				"updatedAt": this.UpdatedAt.Unix(),
				"createdAt": this.CreatedAt.Unix(),
		}
}

func (this *Config) save() error {
		var (
				model = ConfigModelOf()
		)
		this.InitTimes()
		var cnf = model.GetByUnique(this.M())
		if cnf == nil {
				this.InitDefault()
				return model.Add(this)
		}
		(*this.dataClassImpl.data)["updatedAt"] = time.Now().Unix()
		return model.Update(bson.M{"_id": cnf.Id}, this.M(this.formatTime))
}

// 更新时间
func (this *Config) formatTime(m beego.M) beego.M {
		m["createdAt"] = time.Unix(m["createdAt"].(int64), 0)
		m["updatedAt"] = time.Unix(m["updatedAt"].(int64), 0)
		return m
}

func (this *Config) InitTimes() {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now().Local()
		}
}

func (this *Config) setDefaults() {
		if this.State == 0 {
				this.State = 1
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.Root == "" {
				this.Root = DefaultConfigRoot
		}
		this.InitTimes()
}

func (this *Config) setAttributes(data map[string]interface{}, safe ...bool) {
		for key, v := range data {
				if safe[0] {
						// 排除键
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

func (this *Config) Set(key string, v interface{}) *Config {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, v)
		case "key":
				this.SetString(&this.Key, v)
		case "value":
				this.Value = v
		case "root":
				this.SetString(&this.Root, v)
		case "title":
				this.SetString(&this.Title, v)
		case "state":
				this.SetNumInt(&this.State, v)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
		case "updatedAt":
				this.SetTime(&this.UpdatedAt, v)
		}
		return this
}

// 表名
func (this *ConfigModel) TableName() string {
		return ConfigTable
}

// 创建索引
func (this *ConfigModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *ConfigModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"root", "key"},
						Unique: true,
						Sparse: false,
				}))
				this.logs(doc.EnsureIndexKey("state"))
		}
}

// 通过唯一条件查询
func (this *ConfigModel) GetByUnique(data map[string]interface{}, state ...int) *Config {
		var (
				err       error
				cnf       = NewConfig()
				key, root = data["key"], data["root"]
		)
		if key == "" || root == "" {
				return nil
		}
		if len(state) == 0 {
				err = this.FindOne(bson.M{"key": key, "root": root}, cnf)
		} else {
				err = this.FindOne(bson.M{"key": key, "root": root, "state": state[0]}, cnf)
		}
		if err == nil {
				return cnf
		}
		return nil
}

// 设置配置
func (this *ConfigModel) Set(key string, v interface{}, scope ...string) error {
		var data = NewConfig()
		data.Key = key
		data.Value = v
		if len(scope) == 0 {
				scope = append(scope, DefaultConfigRoot)
		}
		data.Root = scope[0]
		data.State = ConfigStateOk
		return data.Save()
}

// 移除配置
func (this *ConfigModel) Remove(key string, scope ...string) error {
		var data = NewConfig()
		data.Key = key
		if len(scope) == 0 {
				scope = append(scope, DefaultConfigRoot)
		}
		data.Root = scope[0]
		data.State = ConfigStateDel
		return data.Save()
}

// 获取对应数据配置
func (this *ConfigModel) GetItemsByScope(scope string, state ...int) []Config {
		if len(state) == 0 {
				state = append(state, ConfigStateOk)
		}
		var (
				err   error
				items []Config
		)
		err = this.Gets(bson.M{"root": scope, "state": state[0]}, &items)
		if err == nil {
				return items
		}
		return items
}

// 获取字符串配置
func (this *ConfigModel) GetString(key string, scope ...string) string {
		if len(scope) == 0 {
				scope = append(scope, DefaultConfigRoot)
		}
		var (
				cnf = NewConfig()
		)
		cnf = this.Get(key, scope[0])
		if cnf == nil {
				return ""
		}
		if cnf.Value == nil {
				return ""
		}
		if str, ok := cnf.Value.(string); ok {
				return str
		}
		return fmt.Sprintf("%v", cnf.Value)
}

// 获取
func (this *ConfigModel) Get(key string, scope string, state ...int) *Config {
		if len(state) == 0 {
				state = append(state, ConfigStateOk)
		}
		var (
				err error
				cnf = NewConfig()
		)
		err = this.FindOne(bson.M{"key": key, "root": scope, "state": ConfigStateOk}, cnf)
		if err == nil {
				return cnf
		}
		return nil
}

// 获取字符串配置
func (this *ConfigModel) GetBool(key string, scope ...string) bool {
		if len(scope) == 0 {
				scope = append(scope, DefaultConfigRoot)
		}
		var (
				cnf = NewConfig()
		)
		cnf = this.Get(key, scope[0])
		if cnf == nil {
				return false
		}
		if cnf.Value == nil {
				return false
		}
		if b, ok := cnf.Value.(bool); ok {
				return b
		}
		// 是否为空
		if IsEmpty(cnf.Value) {
				return false
		}
		// bool
		if str, ok := cnf.Value.(string); ok {
				if strings.EqualFold(str, "true") || strings.EqualFold(str, "yes") {
						return true
				}
				if strings.EqualFold(str, "ok") || strings.EqualFold(str, "on") {
						return true
				}
		}

		return false
}
