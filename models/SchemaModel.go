package models

import (
		"errors"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/transforms"
		"time"
)

// 请求头部
type HeaderSchema struct {
		AppVersion   string  `json:"appVer" bson:"appVer"`                                    //  app 版本 eg: 1.0.0 [0000.00.00]
		Driver       string  `json:"driver" bson:"driver"`                                    //  设备类型 [ios,android,winPhone,win,mac,linux]
		Location     string  `json:"location,omitempty" bson:"location,omitempty"`            //  定位信息 eg: 中国-广东-广州-天河
		Lng          float64 `json:"lng,omitempty" bson:"lng,omitempty"`                      //  定位经度
		Lat          float64 `json:"lat,omitempty" bson:"lat,omitempty"`                      //  定位纬度
		AppId        string  `json:"appId" bson:"appId"`                                      //  appId | android-client:0,ios-client:1,pc:2
		Signature    string  `json:"sign" bson:"signature" `                                  //  参数签名
		Lang         string  `json:"lang" bson:"lang" `                                       //  语言 eg: 中文简体:zh-CN,中文繁体:zh-TW,英语:en, 马来语: my
		Country      string  `json:"country" bson:"country"`                                  //  国家简码 ,eg:CN
		UserOpenId   string  `json:"userOpenId" json:"userOpenId"`                            //  设备 uuid ｜ 用户临时访问身份
		TimeStamp    int64   `json:"timestamp" bson:"timestamp"`                              //  访问时间戳
		Auth         string  `json:"authorization" bson:"authorization"`                      //  用户登陆token
		LaunchLink   string  `json:"launchLink" bson:"launchLink"`                            //  启动页路由｜引导开发链接｜推广链接
		Code         string  `json:"code,omitempty" bson:"code,omitempty"`                    //  追踪码｜活动码
		From         string  `json:"from,omitempty" bson:"from,omitempty"`                    //  来源 0-分享 1-pc 2-h5 3-android 4-ios
		Tags         string  `json:"tags,omitempty" bson:"tags,omitempty"`                    //  tags eg: Travel
		RequestId    string  `json:"requestId,omitempty" bson:"requestId,omitempty"`          //  请求ID
		Crypto       string  `json:"X-Crypto,omitempty" bson:"crypto"`                        //  通信加密方式
		UserAgent    string  `json:"User-Agent" json:"userAgent"`                             //  http client user-agent
		Method       string  `json:":method" bson:"method"`                                   //  请求方法
		Authority    string  `json:":authority,omitempty" json:"authority,omitempty"`         //  请求域
		Path         string  `json:":path,omitempty" bson:"path,omitempty"`                   //  请求接口Path
		ClientIpAddr string  `json:"X-Forwarded-For,omitempty" bson:"clientIpAddr,omitempty"` //  请求ip
		RealIp       string  `json:"X-Real-Ip" bson:"realIp"`                                 //  请求ip
		RemoteAddr   string  `json:"Remote-Address,omitempty" json:"remoteAddr,omitempty"`    //  远程请求ip
		ContentType  string  `json:"Content-Type" bson:"contentType"`                         //  请求数据格式
		Scheme       string  `json:":scheme" bson:"scheme"`                                   //  请求 schema [http,https]
		dataClassImpl
}

// 请求日志
type RequestLog struct {
		HeaderSchema
		Id        bson.ObjectId `json:"id" bson:"_id"`                            // ID
		UserId    string        `json:"userId,omitempty" bson:"userId,omitempty"` // 请求用户ID
		Body      string        `json:"body" bson:"body"`                         // 请求体
		Data      beego.M       `json:"data,omitempty" bson:"data,omitempty"`     // 格式化转换过的数据
		CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`               // 创建时间
}

// 请求日志模型
type RequestLogModel struct {
		BaseModel
}

const (
		RequestLogTableName = "request_logs"
)

var (
		ErrEmptyRequestId = errors.New("empty requestId")
)

func NewHeaderSchema() *HeaderSchema {
		var header = new(HeaderSchema)
		header.Init()
		return header
}

func NewRequestLog() *RequestLog {
		var log = new(RequestLog)
		log.Init()
		return log
}

func RequestLogModelOf() *RequestLogModel {
		var model = new(RequestLogModel)
		model.Init()
		return model
}

func (this *HeaderSchema) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *HeaderSchema) setAttributes(data map[string]interface{}, safe ...bool) {
		for key, v := range data {
				if safe[0] {
						if this.Excludes(key) {
								continue
						}
						if v == nil || v == "" {
								continue
						}
				}
				this.Set(key, v)
		}
}

func (this *HeaderSchema) Set(key string, v interface{}) *HeaderSchema {
		switch key {

		}
		return this
}

func (this *HeaderSchema) setDefaults() {
		if this.Country == "" {
				this.Country = "CN"
		}
		if this.Lang == "" {
				this.Lang = "zh-CN"
		}
		if this.Tags == "" {
				this.Tags = "Travel"
		}
}

func (this *HeaderSchema) data() beego.M {
		return beego.M{
				"appVer":        this.AppVersion,
				"location":      this.Location,
				"lng":           this.Lng,
				"lat":           this.Lat,
				"lang":          this.Lang,
				"appId":         this.AppId,
				"signature":     this.Signature,
				"country":       this.Country,
				"userOpenId":    this.UserOpenId,
				"timestamp":     this.TimeStamp,
				"authorization": this.Auth,
				"launchLink":    this.LaunchLink,
				"code":          this.Code,
				"from":          this.From,
				"tags":          this.Tags,
				"requestId":     this.RequestId,
				"crypto":        this.Crypto,
				"userAgent":     this.UserAgent,
				"method":        this.Method,
				"authority":     this.Authority,
				"path":          this.Path,
				"clientIpAddr":  this.ClientIpAddr,
				"realIp":        this.RealIp,
				"remoteAddr":    this.RemoteAddr,
				"contentType":   this.ContentType,
				"scheme":        this.Scheme,
		}
}

func (this *HeaderSchema) save() error {
		var model = RequestLogModelOf()
		if this.RequestId == "" {
				this.setDefaults()
		}
		if this.RequestId == "" {
				return ErrEmptyRequestId
		}
		var (
				err error
				log = NewRequestLog()
		)
		err = model.GetByKey("requestId", this.RequestId, log)
		if err == nil {
				log.SetAttributes(this.M(), true)
				return model.UpdateById(log.Id.Hex(), log)
		}
		data := this.M()
		data["_id"] = bson.NewObjectId()
		return model.Add(data)
}

func (this *RequestLog) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *RequestLog) setAttributes(data map[string]interface{}, safe ...bool) {
		for key, v := range data {
				if safe[0] {
						if this.Excludes(key) {
								continue
						}
						if v == nil || v == "" {
								continue
						}
				}
				this.Set(key, v)
		}
}

func (this *RequestLog) setDefaults() {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
}

func (this *RequestLog) save() error {
		var model = RequestLogModelOf()
		if this.CreatedAt.IsZero() {
				this.setDefaults()
		}
		if this.RequestId == "" {
				return ErrEmptyRequestId
		}
		var (
				err error
				log = NewRequestLog()
		)
		err = model.GetByKey("requestId", this.RequestId, log)
		if err == nil {
				log.SetAttributes(this.M(), true)
				return model.UpdateById(log.Id.Hex(), log)
		}
		return model.Add(this)
}

func (this *RequestLog) data() beego.M {
		return this.Merger(beego.M{
				"id":        this.Id.Hex(),
				"userId":    this.UserId,
				"data":      this.Data,
				"body":      this.Body,
				"createdAt": this.CreatedAt.Unix(),
		}, this.HeaderSchema.data())
}

func (this *RequestLog) Set(key string, v interface{}) *RequestLog {
		switch key {
		case "id":
				if this.SetObjectId(&this.Id, v) {
						return this
				}
		case "body":
				if this.SetString(&this.Body, v) {
						return this
				}
		case "data":
				if this.SetMapper(&this.Data, v, true) {
						return this
				}
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
				return this
		case "userId":
				this.SetString(&this.UserId, v)
				return this
		}
		this.HeaderSchema.Set(key,v)
		return this
}

// 表名
func (this *RequestLogModel) TableName() string {
		return RequestLogTableName
}

func (this *RequestLogModel) CreateIndex() {
		//userId
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"userId", "clientIpAddr"},
				Unique: false,
				Sparse: true,
		})
		// 唯一索引
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"requestId"},
				Unique: true,
				Sparse: false,
		})
		// 创建索引
		_ = this.Collection().EnsureIndexKey("launchLink")
		_ = this.Collection().EnsureIndexKey("lng", "lat")
		_ = this.Collection().EnsureIndexKey("from", "tags")
		_ = this.Collection().EnsureIndexKey("method", "appId")
		_ = this.Collection().EnsureIndexKey("appVer", "driver")
		_ = this.Collection().EnsureIndexKey("realIp", "remoteAddr")

}
