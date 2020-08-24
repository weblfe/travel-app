package models

import (
		"errors"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/utils"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/transforms"
		"io/ioutil"
		"strings"
		"time"
)

type AddressModel struct {
		BaseModel
}

// 中国行政区域
type Address struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`                // ID
		Level         int           `json:"level" bson:"level"`           // 行政层级
		ParentCode    int64         `json:"parentCode" bson:"parentCode"` // 父级行政代码
		AreaCode      int64         `json:"areaCode" json:"areaCode"`     // 行政代码
		ZipCode       string        `json:"zipCode" bson:"zipCode"`       // 邮政编码
		CityCode      string        `json:"cityCode" bson:"cityCode"`     // 区号
		Name          string        `json:"name" bson:"name"`             // 名称
		ShortName     string        `json:"shortName" bson:"shortName"`   // 简称
		MergerName    string        `json:"mergerName" bson:"mergerName"` // 组合名
		Pinyin        string        `json:"pinyin" bson:"pinyin"`         // 拼音
		Lng           float64       `json:"lng" bson:"lng"`               // 经度
		Lat           float64       `json:"lat" bson:"lat"`               // 纬度
		CreatedAt     time.Time     `json:"createdAt"`                    // 创建时间
		DeletedAt     int64         `json:"deletedAt" json:"deletedAt"`   // 删除时间 ｜ 废弃时间
		dataClassImpl `bson:",omitempty"  json:",omitempty"`
}

const (
		AddressTableName     = "address"
		AddressLevelProvince = 0  // 省
		AddressLevelCity     = 1  // 市
		AddressLevelDistrict = 2  // 区
		AddressLevelTown     = 3  // 镇
		AddressLevelVillage  = 4  // 乡
		AddressLevelUnKnown  = -1 // 异常
)

var (
		_AddressLevelDesc = map[int]string{
				AddressLevelProvince: "省",
				AddressLevelCity:     "市",
				AddressLevelDistrict: "区",
				AddressLevelTown:     "镇",
				AddressLevelVillage:  "乡",
		}

		_AddressLevelMapper = map[int]string{
				AddressLevelProvince: "province",
				AddressLevelCity:     "city",
				AddressLevelDistrict: "district",
				AddressLevelTown:     "town",
				AddressLevelVillage:  "village",
		}
		ErrFileNotExists = errors.New("file not exists")
		ErrEmptyInclude  = errors.New("empty include json")
		ErrEmptyJson     = errors.New("empty json object")
		ErrUnSupportJson = errors.New("unSupport json struct")
)

var (
		MissAreaCodeError = errors.New("miss areaCode")
)

// 新地址
func NewAddress() *Address {
		var addr = new(Address)
		addr.Init()
		return addr
}

// 地址模型
func AddressModelOf() *AddressModel {
		var addr = new(AddressModel)
		addr._Binder = addr
		addr.Init()
		return addr
}

// 文档名
func (this *AddressModel) TableName() string {
		return AddressTableName
}

// 创建索引
func (this *AddressModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *AddressModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				// unique areaCode
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"areaCode"},
						Unique: true,
						Sparse: false,
				}))
				// index
				this.logs(doc.EnsureIndexKey("level", "deletedAt"))
				this.logs(doc.EnsureIndexKey("name", "shortName", "pinyin"))
				this.logs(doc.EnsureIndexKey("lng", "lat"))
		}
}

// 地址查询
func (this *AddressModel) GetAddress(query map[string]interface{}) *Address {
		if query == nil || len(query) == 0 {
				return nil
		}
		var addr = NewAddress()
		table := this.Collection()
		defer this.destroy(table)
		if err := table.Find(beego.M(query)).One(addr); err == nil {
				return addr
		}
		return nil
}

// 通过json 倒入
func (this *AddressModel) ImportFromJsonFile(file string) (int, error) {
		var data = beego.M{}
		if !utils.FileExists(file) {
				return 0, ErrFileNotExists
		}
		jsonByte, err := ioutil.ReadFile(file)
		if err != nil {
				return 0, err
		}
		err = libs.Json().Unmarshal(jsonByte, &data)
		if err != nil {
				return 0, err
		}
		if len(data) == 0 {
				return 0, ErrEmptyJson
		}
		var items, ok = data["RECORDS"]
		if !ok {
				return 0, ErrUnSupportJson
		}
		var (
				arr  []map[string]interface{}
				docs = make([]interface{}, len(arr))
		)
		docs = docs[:0]
		arr, ok = items.([]map[string]interface{})
		if ok && len(arr) > 0 {
				cache := map[int64]bool{}
				for _, it := range arr {
						if len(it) <= 0 {
								continue
						}
						it = this.transformJson(it)
						v, ok := it["areaCode"]
						if !ok || v == "" || v == nil || v == 0 {
								continue
						}
						if _, ok := cache[v.(int64)]; ok {
								continue
						}
						cache[v.(int64)] = true
						if this.Exists(beego.M{"areaCode": v}) {
								continue
						}
						docs = append(docs, it)
				}
		}
		if len(docs) == 0 {
				return 0, ErrEmptyInclude
		}
		err = this.Inserts(docs)
		return len(docs), err
}

// 转换
func (this *AddressModel) transformJson(data map[string]interface{}) map[string]interface{} {
		var transMapper = map[string]string{
				"area_code":   "areaCode",
				"parent_code": "parentCode",
				"zip_code":    "zipCode",
				"city_code":   "cityCode",
				"short_name":  "shortName",
				"merger_name": "mergerName",
		}
		for key, NewKey := range transMapper {
				if v, ok := data[key]; ok {
						data[NewKey] = v
						delete(data, key)
				}
		}
		return data
}

// 地址对象初始化
func (this *Address) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

// 数据输出
func (this *Address) data() beego.M {
		return beego.M{
				"id":         this.Id.Hex(),
				"level":      this.Level,
				"levelCode":  GetLevelText(this.Level),
				"levelDesc":  GetLevelDesc(this.Level),
				"parentCode": this.ParentCode,
				"areaCode":   this.AreaCode,
				"zipCode":    this.ZipCode,
				"cityCode":   this.CityCode,
				"name":       this.Name,
				"shortName":  this.ShortName,
				"mergerName": this.MergerName,
				"pinyin":     this.Pinyin,
				"lng":        this.Lng,
				"lat":        this.Lat,
				"createdAt":  this.CreatedAt,
				"deletedAt":  this.DeletedAt,
		}
}

// 加载数据
func (this *Address) Load(data map[string]interface{}) *Address {
		this.setAttributes(data)
		return this
}

// 设置数据值
func (this *Address) setAttributes(data map[string]interface{}, safe ...bool) {
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

// setter
func (this *Address) Set(key string, v interface{}) *Address {
		switch key {
		case "id":
				if id, ok := v.(bson.ObjectId); ok && this.Id == "" {
						this.Id = id
						return this
				}
				if id, ok := v.(*bson.ObjectId); ok && this.Id == "" {
						this.Id = *id
						return this
				}
				if id, ok := v.(string); ok && this.Id == "" {
						this.Id = bson.ObjectIdHex(id)
						return this
				}
		case "level":
				if n, ok := v.(int); ok && InSupportLevel(n) {
						this.Level = n
						return this
				}
				if n, ok := v.(string); ok {
						l := GetLevelByDesc(n)
						if !InSupportLevel(l) {
								return this
						}
						this.Level = l
						return this
				}
		case "parentCode":
				this.ParentCode = v.(int64)
		case "areaCode":
				this.AreaCode = v.(int64)
		case "zipCode":
				if code, ok := v.(int); ok && code > 0 {
						this.ZipCode = fmt.Sprintf("%d", code)
						return this
				}
				this.ZipCode = v.(string)
		case "cityCode":
				if code, ok := v.(int); ok && code > 0 {
						this.CityCode = fmt.Sprintf("%d", code)
						return this
				}
				this.CityCode = v.(string)
		case "name":
				this.Name = v.(string)
		case "shortName":
				this.ShortName = v.(string)
		case "mergerName":
				if arr, ok := v.([]string); ok {
						this.MergerName = strings.Join(arr, ",")
						return this
				}
				this.MergerName = v.(string)
		case "pinyin":
				this.Pinyin = v.(string)
		case "lng":
				this.Lng = v.(float64)
		case "lat":
				this.Lat = v.(float64)
		case "createdAt":
				if t, ok := v.(time.Time); ok {
						this.CreatedAt = t
						return this
				}
				if t, ok := v.(*time.Time); ok {
						this.CreatedAt = *t
						return this
				}
				if t, ok := v.(int64); ok && t > 0 {
						this.CreatedAt = time.Unix(t, 0)
				}
		case "deletedAt":
				if t, ok := v.(time.Time); ok {
						this.DeletedAt = t.Unix()
						return this
				}
				if t, ok := v.(*time.Time); ok {
						this.DeletedAt = t.Unix()
						return this
				}
				if t, ok := v.(int64); ok && t > 0 {
						this.DeletedAt = t
				}
		}
		return this
}

// 设置默认值
func (this *Address) setDefaults() {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
}

// 默认值填充
func (this *Address) Defaults() *Address {
		this.setDefaults()
		return this
}

// 保存器
func (this *Address) save() error {
		var model = AddressModelOf()
		if this.AreaCode == 0 {
				return MissAreaCodeError
		}
		var addr = NewAddress()
		err := model.GetByKey("areaCode", this.AreaCode, addr)
		if err != nil && model.IsNotFound(err) {
				return model.Add(this)
		}
		this.SetAttributes(addr.M(), true)
		return model.UpdateById(addr.Id.Hex(), this.M())
}

// 获取行政区域 文本
func GetLevelText(level int) string {
		return _AddressLevelMapper[level]
}

// 获取行政区域 描述
func GetLevelDesc(level int) string {
		return _AddressLevelDesc[level]
}

// 是否支持的行政等级
func InSupportLevel(level int) bool {
		for l, _ := range _AddressLevelDesc {
				if l == level {
						return true
				}
		}
		return false
}

// 通过表现获取行政等级
func GetLevelByDesc(desc string) int {
		var descTrim = strings.TrimSpace(desc)
		for level, d := range _AddressLevelDesc {
				if desc == d || descTrim == d {
						return level
				}
		}
		return AddressLevelUnKnown
}
