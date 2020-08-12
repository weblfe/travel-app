package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/siddontang/go/bson"
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
		Extras                bson.M        `json:"extras" bson:"extras"`                               // 备用信息
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

func (this *PopularizationChannels) Init() {
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.defaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *PopularizationChannels) Set(key string, v interface{}) *PopularizationChannels {

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

func PopularizationChannelsModelOf() *PopularizationChannelsModel {
		var model = new(PopularizationChannelsModel)
		model._Self = model
		model.Init()
		return model
}

func (this *PopularizationChannelsModel) TableName() string {
		return PopularizationChannelsTable
}

func (this *PopularizationChannelsModel) CreateIndex() {
		// unique mobile
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"channel"},
				Unique: true,
				Sparse: false,
		})
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"wechat"},
				Unique: true,
				Sparse: false,
		})
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"userId", "name"},
				Unique: true,
				Sparse: false,
		})
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"mobile", "name"},
				Unique: true,
				Sparse: false,
		})
		_ = this.Collection().EnsureIndexKey("email")
}

func (this *PopularizationChannelsModel) GetByUnique(m beego.M) *PopularizationChannels {
		return nil
}
