package models

import (
		"errors"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"time"
)

type PopularizationChannelsModel struct {
		BaseModel
}

// 奖励方式 ｜ 结算方式
type AwardWay struct {
		Name   string `json:"name" bson:"name"`
		Config bson.M `json:"config" bson:"config"`
}

// 权限
type Permission struct {
		Name string `json:"name" bson:"name"`
		On   int    `json:"on" bson:"on"`
}

// 推广渠道记录
type PopularizationChannels struct {
		Id                    bson.ObjectId `json:"id" bson:"_id"`                                      // 渠道ID
		Name                  string        `json:"name" bson:"name"`                                   // 渠道名
		UserId                string        `json:"userId,omitempty" bson:"userId,omitempty"`           // 用户ID
		Mobile                string        `json:"mobile,omitempty" bson:"mobile,omitempty"`           // 推广用户｜渠道联系号码
		QrcodeUrl             string        `json:"qrcodeUrl" bson:"qrcodeUrl"`                         // 二维码链接
		Email                 string        `json:"email,omitempty" bson:"email,omitempty"`             // 渠道邮箱
		WeChat                string        `json:"wechat" bson:"wechat"`                               // 微信号
		Awards                []*AwardWay   `json:"awards" bson:"awards"`                               // 奖励方式
		InvitedRegisterNumber int64         `json:"invitedRegisterNumber" bson:"invitedRegisterNumber"` // 邀请注册人数
		ValidNumber           int64         `json:"validNumber" bson:"validNumber"`                     // 有效邀请注册人数
		Channel               string        `json:"channel" bson:"channel"`                             // 6-128 渠道码
		Extras                beego.M       `json:"extras" bson:"extras"`                               // 备用信息
		Status                int           `json:"status" bson:"status"`                               // 状态 ： 0 ，1， 2
		ParentId              bson.ObjectId `json:"parentId,omitempty" bson:"parentId,omitempty"`       // 父级推广渠道
		Permissions           []*Permission `json:"permissions" bson:"permissions"`                     // 渠道权限
		Comment               string        `json:"comment" bson:"comment"`                             // 备注
		CreatedAt             time.Time     `json:"createdAt" bson:"createdAt"`                         // 创建时间
		UpdatedAt             time.Time     `json:"updatedAt" bson:"updatedAt"`                         // 更新时间
		dataClassImpl         `bson:",omitempty"  json:",omitempty"`
}

const (
		PopularizationChannelsTable = "popularization_channels"
)

func NewPopularizationChannel() *PopularizationChannels {
		var channelInfo = new(PopularizationChannels)
		channelInfo.Init()
		return channelInfo
}

func (this *PopularizationChannels) Verify() error {
		if this.Channel == "" {
				return errors.New("miss channel ")
		}
		if this.Name == "" {
				return errors.New("miss name of channel")
		}
		if this.Mobile == "" && this.UserId == "" && this.WeChat == "" {
				return errors.New("miss channel user")
		}
		if this.ParentId != "" && !PopularizationChannelsModelOf().Exists(bson.M{"_id": this.ParentId}) {
				return errors.New("parent channel not exists")
		}
		return nil
}

func (this *PopularizationChannels) Incr(key string, incr ...int) bool {
		if len(incr) == 0 {
				incr = append(incr, 1)
		}
		if key == "validNumber" {
				this.ValidNumber = this.ValidNumber + int64(incr[0])
				return true
		}
		if key == "invitedRegisterNumber" {
				this.InvitedRegisterNumber = this.InvitedRegisterNumber + int64(incr[0])
				return true
		}
		return false
}

func (this *PopularizationChannels) Init() {
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.defaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *PopularizationChannels) Set(key string, v interface{}) *PopularizationChannels {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, v)
		case "name":
				this.SetString(&this.Name, v)
		case "userId":
				this.SetString(&this.UserId, v)
		case "mobile":
				this.SetString(&this.Mobile, v)
		case "qrcodeUrl":
				this.SetString(&this.QrcodeUrl, v)
		case "email":
				this.SetString(&this.Email, v)
		case "wechat":
				this.SetString(&this.WeChat, v)
		case "awards":
				if arr, ok := v.([]*AwardWay); ok {
						this.Awards = arr
				}
		case "invitedRegisterNumber":
				this.SetNumIntN(&this.InvitedRegisterNumber, v)
		case "validNumber":
				this.SetNumIntN(&this.ValidNumber, v)
		case "channel":
				this.SetString(&this.Channel, v)
		case "extras":
				this.SetMapper(&this.Extras, v)
		case "parentId":
				this.SetObjectId(&this.ParentId, v)
		case "permissions":
		case "comment":
				this.SetString(&this.Comment, v)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
		case "updatedAt":
				this.SetTime(&this.UpdatedAt, v)
		}
		return this
}

func (this *PopularizationChannels) data() beego.M {
		return beego.M{
				"id":                    this.Id.Hex(),
				"userId":                this.UserId,
				"name":                  this.Name,
				"mobile":                this.getMobile(),
				"qrcodeUrl":             this.getQrcodeUrl(),
				"email":                 this.getEmail(),
				"weChat":                this.getWeChat(),
				"awards":                this.Awards,
				"invitedRegisterNumber": this.InvitedRegisterNumber,
				"validNumber":           this.ValidNumber,
				"channel":               this.Channel,
				"extras":                this.Extras,
				"status":                this.Status,
				"parentId":              this.ParentId,
				"permissions":           this.Permissions,
				"comment":               this.Comment,
				"updatedAt":             this.UpdatedAt.Unix(),
				"createdAt":             this.CreatedAt.Unix(),
		}
}

func (this *PopularizationChannels) getMobile() string {
		if this.Mobile != "" {
				return this.Mobile
		}
		return ""
}

func (this *PopularizationChannels) getQrcodeUrl() string {
		if this.QrcodeUrl != "" {
				return this.QrcodeUrl
		}
		return ""
}

func (this *PopularizationChannels) getEmail() string {
		if this.Email != "" {
				return this.Email
		}
		return ""
}

func (this *PopularizationChannels) getWeChat() string {
		if this.WeChat != "" {
				return this.WeChat
		}
		return ""
}

func (this *PopularizationChannels) setAttributes(data beego.M, safe ...bool) {
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

func (this *PopularizationChannels) save() error {
		var (
				tmp   = NewPopularizationChannel()
				model = PopularizationChannelsModelOf()
		)
		if this.Id != "" {
				err := model.GetById(this.Id.Hex(), tmp)
				if err == nil {
						return model.Update(beego.M{"_id": this.Id}, this.M())
				}
		}
		tmp = model.GetByUnique(this.data())
		if tmp != nil {
				return model.Update(beego.M{"_id": tmp.Id}, this.M())
		}
		this.reset()
		this.InitDefault()
		return model.Add(this)
}

func (this *PopularizationChannels) defaults() {
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now().Local()
		}
		if this.Awards == nil {
				this.Awards = make([]*AwardWay, 2)
				this.Awards = this.Awards[:0]
		}
		if this.Permissions == nil {
				this.Permissions = make([]*Permission, 2)
				this.Permissions = this.Permissions[:0]
		}
		if this.Status == 0 {
				this.Status = 1
		}
}

func (this *PopularizationChannels) IsOk() bool {
		if this.Status == StatusOk {
				return true
		}
		return false
}

func PopularizationChannelsModelOf() *PopularizationChannelsModel {
		var model = new(PopularizationChannelsModel)
		model.Bind(model)
		model.Init()
		return model
}

func (this *PopularizationChannelsModel) TableName() string {
		return PopularizationChannelsTable
}

func (this *PopularizationChannelsModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *PopularizationChannelsModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"channel"},
						Unique: true,
						Sparse: false,
				}))
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"wechat"},
						Unique: true,
						Sparse: false,
				}))
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"userId", "name"},
						Unique: true,
						Sparse: false,
				}))
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"mobile", "name"},
						Unique: true,
						Sparse: false,
				}))
				this.logs(doc.EnsureIndexKey("email"))
		}
}

func (this *PopularizationChannelsModel) GetByUnique(m beego.M) *PopularizationChannels {
		var (
				err    error
				query  beego.M
				object *PopularizationChannels
				keys   = [][]string{
						{"channel"},
						{"wechat"},
						{"userId", "name"},
						{"mobile", "name"},
				}
		)
		if len(m) == 0 {
				return nil
		}
		for _, arr := range keys {
				query = make(beego.M)
				for _, v := range arr {
						value, ok := m[v]
						if !ok {
								query = nil
								break
						}
						query[v] = value
				}
				if query == nil || len(query) == 0 {
						continue
				}
				object = NewPopularizationChannel()
				err = this.FindOne(query, object)
				if err == nil && object != nil {
						return object
				}
		}
		return nil
}
