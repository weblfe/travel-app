package models

import (
		"crypto/md5"
		"encoding/hex"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/transforms"
		"strings"
		"time"
)

type UserAddressModel struct {
		BaseModel
}

// 地址模型
type UserAddress struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`                                    // id
		UserId        bson.ObjectId `json:"userId" bson:"userId"`                             // 地址所属用户
		Type          int           `json:"type" bson:"type"`                                 // 地址类型 1: 注册地址
		Version       int           `json:"version" json:"version"`                           // 地址更新处理
		Hash          string        `json:"hash" bson:"hash"`                                 // 用户地址唯一值
		Sort          int           `json:"sort" json:"sort"`                                 // 排序
		Text          string        `json:"text,omitempty" json:"text,omitempty"`             // 字符串文本地址
		State         int           `json:"state" bson:"state"`                               // 状态 -1:已过期,0:不可用, 1: 可用｜正常
		Country       string        `json:"country,omitempty" bson:"country,omitempty"`       // 国家
		CountryId     int           `json:"countryId" bson:"countryId"`                       // 国家编码
		City          string        `json:"city,omitempty" bson:"city,omitempty"`             // 城市
		CityId        string        `json:"cityId,omitempty" bson:"cityId,omitempty"`         // 城市编码
		Province      string        `json:"province,omitempty" bson:"province,omitempty"`     // 省份
		ProvinceId    string        `json:"provinceId,omitempty" bson:"provinceId,omitempty"` // 省份编码
		District      string        `json:"district,omitempty" bson:"district,omitempty"`     // 行政区
		DistrictId    string        `json:"districtId,omitempty" json:"districtId,omitempty"` // 行政区编码
		Street        string        `json:"street,omitempty" json:"street,omitempty"`         // 街区
		StreetId      string        `json:"streetId,omitempty" json:"streetId,omitempty"`     // 街区编码
		Floor         string        `json:"floor,omitempty" bson:"floor,omitempty"`           // 楼层  eg : 天辉大厦18楼
		Doorplate     string        `json:"doorplate,omitempty" bson:"doorplate,omitempty"`   // 门牌号
		Longitude     float64       `json:"longitude,omitempty" bson:"longitude,omitempty"`   // 经度 [东-西]
		Latitude      float64       `json:"latitude,omitempty" bson:"latitude,omitempty"`     // 纬度 [南-北]
		CreateTime    int64         `json:"createTime" bson:"createTime"`                     // 创建日期
		dataClassImpl `bson:",omitempty"  json:",omitempty"`
}

const (
		AddressTypeTmp        = 0
		AddressTypeRegister   = 1
		AddressTypeCompany    = 2
		AddressTypeHome       = 3
		AddressTypeSchool     = 4
		AddressTypeLogin      = 5
		AddressTypeVirtual    = 6
		AddressTypeExpress    = 7
		AddressStateNil       = 0
		AddressStateActivate  = 1
		AddressStateDelete    = 2
		AddressStateExpired   = 3
		UserAddressModelTable = "user_address"
		DefaultCountry        = "中国"
		DefaultCountryId      = 48
)

var (
		// 类型描述
		_AddressTypeDesc = map[int]string{
				AddressTypeTmp:      "临时地址",
				AddressTypeRegister: "注册地址",
				AddressTypeCompany:  "公司",
				AddressTypeHome:     "家",
				AddressTypeSchool:   "学校",
				AddressTypeLogin:    "登陆",
				AddressTypeVirtual:  "虚拟",
				AddressTypeExpress:  "快递",
		}
		// 类型名
		_AddressTypeMapper = map[int]string{
				AddressTypeTmp:      "tmp",
				AddressTypeRegister: "register",
				AddressTypeCompany:  "company",
				AddressTypeHome:     "home",
				AddressTypeSchool:   "school",
				AddressTypeLogin:    "login",
				AddressTypeVirtual:  "virtual",
				AddressTypeExpress:  "express",
		}
		// 地址状态描述
		_AddressStateDesc = map[int]string{
				AddressStateNil:      "未激活",
				AddressStateActivate: "正常",
				AddressStateDelete:   "已删除",
				AddressStateExpired:  "已过期",
		}
)

func NewUserAddress() *UserAddress {
		var address = new(UserAddress)
		address.Init()
		return address
}

func GetStateDesc(state int) string {
		return _AddressStateDesc[state]
}

func GetTypeCode(typ int) string {
		return _AddressTypeMapper[typ]
}

func GetTypeId(typ string) int {
		for id, k := range _AddressTypeMapper {
				if k == typ || strings.EqualFold(k, typ) {
						return id
				}
		}
		return AddressTypeTmp
}

func GetTypeIdByDesc(desc string) int {
		var descTrim = strings.TrimSpace(desc)
		for id, k := range _AddressTypeDesc {
				if k == desc || descTrim == k {
						return id
				}
		}
		return AddressTypeTmp
}

func GetTypeDesc(typ int) string {
		return _AddressTypeDesc[typ]
}

func UserAddressModelOf() *UserAddressModel {
		var model = new(UserAddressModel)
		model._Binder = model
		model.Init()
		return model
}

func (this *UserAddress) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *UserAddress) Parse(v interface{}) bool {
		switch v.(type) {
		case UserAddress:
				addr := v.(UserAddress)
				this.SetAttributes(addr.M(), true)
				this.InitDefault()
				return true
		case *UserAddress:
				addr := v.(*UserAddress)
				this.SetAttributes(addr.M(), true)
				this.InitDefault()
				return true
		case string:
				if this.Text == "" && v != "" {
						this.Text = v.(string)
				}
				ok := this.setByText()
				this.InitDefault()
				return ok
		}
		return false
}

func (this *UserAddress) setByText() bool {
		if this.Text == "" {
				return false
		}
		var addrArr = strings.SplitN(this.Text, "-", -1)
		if len(addrArr) == 1 {
				addrArr = strings.SplitN(this.Text, " ", -1)
		}
		var addrLen = len(addrArr)
		if addrLen < 3 {
				return false
		}
		this.Province = addrArr[0] // 广东
		this.City = addrArr[1]     // 广州
		this.District = addrArr[2] // 越秀
		if addrLen >= 4 {
				this.Doorplate = addrArr[3] // xx路xxx
		}
		// @todo
		/*if strings.Contains(this.Text,"@lng:") {
		}
		if strings.Contains(this.Text,"@lat:") {
		}*/
		return true
}

// 数据加载
func (this *UserAddress) Load(data map[string]interface{}) *UserAddress {
		this.setAttributes(data)
		return this
}

// 设置器
func (this *UserAddress) Set(key string, v interface{}) *UserAddress {
		switch key {
		case "userId":
				if id, ok := v.(bson.ObjectId); ok {
						this.UserId = id
						return this
				}
				if id, ok := v.(*bson.ObjectId); ok {
						this.UserId = *id
						return this
				}
				if id, ok := v.(string); ok {
						this.UserId = bson.ObjectIdHex(id)
						return this
				}
		case "type":
				this.Type = v.(int)
		case "sort":
				this.Sort = v.(int)
		case "text":
				this.Text = v.(string)
		case "state":
				this.State = v.(int)
		case "country":
				this.Country = v.(string)
		case "countryId":
				this.CountryId = v.(int)
		case "city":
				this.City = v.(string)
		case "cityId":
				if str, ok := v.(string); ok {
						this.CityId = str
						return this
				}
				if id, ok := v.(int); ok {
						this.CityId = fmt.Sprintf("%d", id)
				}
		case "province":
				this.Province = v.(string)
		case "provinceId":
				if str, ok := v.(string); ok {
						this.ProvinceId = str
						return this
				}
				if id, ok := v.(int); ok {
						this.ProvinceId = fmt.Sprintf("%d", id)
				}
		case "district":
				this.District = v.(string)
		case "districtId":
				if str, ok := v.(string); ok {
						this.DistrictId = str
						return this
				}
				if id, ok := v.(int); ok {
						this.DistrictId = fmt.Sprintf("%d", id)
				}
		case "street":
				this.Street = v.(string)
		case "streetId":
				if str, ok := v.(string); ok {
						this.StreetId = str
						return this
				}
				if id, ok := v.(int); ok {
						this.StreetId = fmt.Sprintf("%d", id)
				}
		case "floor":
				this.Floor = v.(string)
		case "doorplate":
				this.Doorplate = v.(string)
		case "longitude":
				this.Longitude = v.(float64)
		case "latitude":
				this.Latitude = v.(float64)
		case "createTime":
				if t, ok := v.(time.Time); ok {
						this.CreateTime = t.Unix()
						return this
				}
				if n, ok := v.(int64); ok {
						this.CreateTime = n
						return this
				}
				if n, ok := v.(int); ok {
						this.CreateTime = int64(n)
						return this
				}
		}
		return this
}

// 设置数据属性
func (this *UserAddress) setAttributes(data map[string]interface{}, safe ...bool) {
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

// 设置默认值
func (this *UserAddress) setDefaults() {
		if this.CreateTime == 0 {
				this.CreateTime = time.Now().Unix()
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.State == AddressStateNil {
				this.State = AddressStateActivate
		}
		if this.Hash == "" {
				this.Hash = this.getHash()
		}
		if this.Sort == 0 {
				this.Sort = 1
		}
		if this.Version == 0 {
				this.Version = 1
		}
		if this.Text == "" {
				this.Text = this.String()
		}
		if this.Country == "" {
				this.Country = DefaultCountry
		}
}

// 获取hash 值
func (this *UserAddress) getHash() string {
		var ins = md5.New()
		ins.Write([]byte(this.toString()))
		return hex.EncodeToString(ins.Sum(nil))
}

// 地址
func (this *UserAddress) toString() string {
		var hashTemp = `userId=%s&type=%d&version=%d&sort=%d`
		hashTemp += `&text=%s&country=%s&countryId=%d`
		hashTemp += `&province=%s&provinceId=%d&city=%s&`
		hashTemp += `&cityId=%s&district=%s&districtId=%s&`
		hashTemp += `&street=%s&streetId=%s&floor=%s&doorplate=%s`
		hashTemp += `&longitude=%f&latitude=%f&createTime=%d`
		return fmt.Sprintf(hashTemp,
				this.UserId.Hex(), this.Type, this.Version, this.Sort, this.Text, this.Country,
				this.CountryId, this.Province, this.ProvinceId, this.City, this.CityId,
				this.District, this.DistrictId, this.Street, this.StreetId, this.Floor, this.Doorplate,
				this.Longitude, this.Latitude, this.CreateTime,
		)
}

// 更新地址版本号
func (this *UserAddress) IncrVersion() *UserAddress {
		this.Version++
		return this
}

// 更新数据数据逻辑
func (this *UserAddress) Update() error {
		this.IncrVersion()
		this.Hash = this.getHash()
		return this.save()
}

func (this *UserAddress) data() beego.M {
		return beego.M{
				"id":         this.Id.Hex(),
				"userId":     this.UserId.Hex(),
				"type":       this.Type,
				"typeDesc":   GetTypeDesc(this.Type),
				"sort":       this.Sort,
				"state":      this.State,
				"stateDesc":  GetStateDesc(this.State),
				"hash":       this.Hash,
				"text":       this.Text,
				"country":    this.getCountry(),
				"countryId":  this.getCountryId(),
				"province":   this.Province,
				"provinceId": this.ProvinceId,
				"city":       this.City,
				"cityId":     this.CityId,
				"district":   this.District,
				"districtId": this.DistrictId,
				"street":     this.Street,
				"streetId":   this.StreetId,
				"floor":      this.Floor,
				"doorplate":  this.Doorplate,
				"latitude":   this.Latitude,
				"longitude":  this.Longitude,
				"address":    this.String(),
				"version":    this.Version,
				"createTime": this.CreateTime,
		}
}

func (this *UserAddress) getCountry() string {
		if this.Country == "" {
				return DefaultCountry
		}
		return this.Country
}

func (this *UserAddress) getCountryId() int {
		if this.CountryId == 0 {
				return DefaultCountryId
		}
		return this.CountryId
}

func (this *UserAddress) String() string {
		if this.Text != "" {
				return this.Text
		}
		var text = fmt.Sprintf("%s-%s-%s-%s-%s", this.Country, this.Province, this.City, this.District, this.Doorplate)
		if text == "----" {
				return ""
		}
		return text
}

func (this *UserAddress) Location() *Location {
		var location = new(Location)
		location.Latitude = this.Latitude
		location.Longitude = this.Longitude
		return location
}

func (this *UserAddress) save() error {
		var (
				model = UserAddressModelOf()
				hash  = this.getHash()
		)
		if this.Id != "" && model.Exists(bson.M{"_id": this.Id}) {
				return model.Update(beego.M{"_id": this.Id}, this)
		}
		var addr = NewUserAddress()
		if this.Hash == "" || this.Hash != hash {
				this.Hash = hash
				this.InitDefault()
				return model.Add(this)
		}
		err := model.GetByKey("hash", hash, addr)
		if err == nil {
				return model.Update(beego.M{"_id": this.Id}, this)
		}
		return err
}

func (this *UserAddressModel) TableName() string {
		return UserAddressModelTable
}

func (this *UserAddressModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *UserAddressModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"hash"},
						Unique: true,
						Sparse: false,
				}))
				// 创建索引
				this.logs(doc.EnsureIndexKey("doorplate"))
				this.logs(doc.EnsureIndexKey("city", "district"))
				this.logs(doc.EnsureIndexKey("country", "province"))
				this.logs(doc.EnsureIndexKey("longitude", "latitude"))
				this.logs(doc.EnsureIndexKey("userId", "type", "state"))
				this.logs(doc.EnsureIndexKey("countryId", "provinceId", "cityId", "districtId"))
		}
}

func (this *UserAddressModel) GetUserAddress(userId bson.ObjectId, typ int) string {
		var addr = this.GetAddressByUserId(userId, typ)
		if addr != nil {
				return addr.String()
		}
		return ""
}

func (this *UserAddressModel) GetAddressByUserId(userId bson.ObjectId, typ ...int) *UserAddress {
		table := this.Document()
		defer this.destroy(table)
		var (
				addr  = new(UserAddress)
				query = bson.M{"userId": userId, "state": AddressStateActivate}
		)
		if len(typ) != 0 {
				query["type"] = typ[0]
		}
		err := table.Find(query).Sort("+sort").One(addr)
		if err == nil {
				return addr
		}
		return nil
}

func (this *UserAddressModel) GetUserAddressList(userId bson.ObjectId, typ int) []string {
		table := this.Document()
		defer this.destroy(table)
		var (
				addrArr = make([]*UserAddress, 2)
				arr     []string
				txt     string
		)
		addrArr = addrArr[:0]
		err := table.Find(bson.M{"userId": userId, "type": typ, "state": AddressStateActivate}).Sort("+sort").All(&addrArr)
		if err == nil {
				for _, addr := range addrArr {
						if addr == nil {
								continue
						}
						txt = addr.String()
						if txt == "" {
								continue
						}
						arr = append(arr, txt)
				}
		}
		return arr
}

// 定位
type Location struct {
		Longitude float64 `json:"longitude" bson:"longitude"` // 经度 [东-西]
		Latitude  float64 `json:"latitude" bson:"latitude"`   // 纬度 [南-北]
}

// 3d 定位
type Location3d struct {
		Location
		Altitude float64 `json:"altitude" bson:"altitude"` // 高度 [海拔]
}

func NewLocation(x, y float64) *Location {
		var location = new(Location)
		location.Longitude = x
		location.Latitude = y
		return location
}

func (this *Location) PointX() float64 {
		return this.Longitude
}

func (this *Location) PointY() float64 {
		return this.Latitude
}

func (this *Location) Points() []float64 {
		return []float64{this.PointX(), this.PointY()}
}

func (this *Location) GetLongitudeDesc() string {
		if this.Longitude < 0 {
				return fmt.Sprintf("西经:%f", this.Longitude)
		}
		return fmt.Sprintf("东经:%f", this.Longitude)
}

func (this *Location) GetLatitudeDesc() string {
		if this.Latitude < 0 {
				return fmt.Sprintf("南纬:%f", this.Latitude)
		}
		return fmt.Sprintf("北纬:%f", this.Latitude)
}

func (this *Location) GetDesc() []string {
		return []string{this.GetLongitudeDesc(), this.GetLatitudeDesc()}
}

func (this *Location) Json() []byte {
		data, _ := libs.Json().Marshal(this)
		return data
}

func NewLocation3d(x, y, z float64) *Location3d {
		var location = new(Location3d)
		location.Longitude = x
		location.Latitude = y
		location.Altitude = z
		return location
}

func (this *Location3d) PointZ() float64 {
		return this.Altitude
}

func (this *Location3d) GetAltitudeDesc() string {
		return fmt.Sprintf("海拔:%f", this.Altitude)
}

func (this *Location3d) GetDesc() []string {
		return []string{this.GetLongitudeDesc(), this.GetLatitudeDesc(), this.GetAltitudeDesc()}
}

func (this *Location3d) Points() []float64 {
		return []float64{this.PointX(), this.PointY(), this.PointZ()}
}

func (this *Location3d) Json() []byte {
		data, _ := libs.Json().Marshal(this)
		return data
}
