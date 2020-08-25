package services

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
		"time"
)

type PopularizationService interface {
		GetChannelQrcode(ch string) string
		GetChannelInfo(ch string) *models.PopularizationChannels
}

type popularizationServiceImpl struct {
		BaseService
		model *models.PopularizationChannelsModel
}

const (
		PopularizationServiceClass = "PopularizationService"
)

func PopularizationServiceOf() *popularizationServiceImpl {
		var service = new(popularizationServiceImpl)
		service.Init()
		return service
}

func (this *popularizationServiceImpl) Init() {
		this.ClassName = PopularizationServiceClass
		this.Constructor = func(args ...interface{}) interface{} {
				return PopularizationServiceOf()
		}
		this.model = models.PopularizationChannelsModelOf()
		this.init()
}

func (this *popularizationServiceImpl) GetChannelInfo(ch string) *models.PopularizationChannels {
		var (
				data = models.NewPopularizationChannel()
				err  = this.model.NewQuery(bson.M{"channel": ch}).One(data)
		)
		if err != nil {
				return nil
		}
		if data.QrcodeUrl == "" {
				this.createQrcode(data)
		}
		return data
}

func (this *popularizationServiceImpl) GetChannelQrcode(ch string) string {
		var info = this.GetChannelInfo(ch)
		if info == nil {
				return ""
		}
		if info.QrcodeUrl == "" {
				this.createQrcode(info)
		}
		return info.QrcodeUrl
}

func (this *popularizationServiceImpl) createQrcode(data *models.PopularizationChannels) string {
		var (
				url    string
				err    error
				params = NewQrcodeParams()
		)
		params.Init()
		params.Set("userId", data.UserId).Set("channel", data.Channel).Set("referId",data.Id.Hex())
		url, err = QrcodeServiceOf().CreateQrcode(params)
		if err == nil {
				return url
		}
		return ""
}

func (this *popularizationServiceImpl) Exists(ch string, status ...int) bool {
		var (
				n   int
				err error
		)
		if len(status) <= 0 {
				n, err = this.model.NewQuery(bson.M{"channel": ch}).Count()
		} else {
				n, err = this.model.NewQuery(bson.M{"channel": ch, "status": status[0]}).Count()
		}
		if err != nil || n == 0 {
				return false
		}
		return true
}

func (this *popularizationServiceImpl) Update(channel string, data beego.M) error {
		if len(data) == 0 {
				return common.NewErrors(common.EmptyParamCode, "更新参数不能为空")
		}
		data["updatedAt"] = time.Now().Local()
		return this.model.Update(bson.M{"channel": channel}, data)
}

func (this *popularizationServiceImpl) Check(channels *models.PopularizationChannels) error {
		return channels.Verify()
}

func (this *popularizationServiceImpl) Create(channels *models.PopularizationChannels) error {
		var err error
		if channels == nil {
				return common.NewErrors(common.EmptyParamCode, "创建异常")
		}
		channels.InitDefault()
		err = this.Check(channels)
		if err != nil {
				return err
		}
		if channels.Status == models.StatusOff {
				return common.NewErrors(common.ErrorCode, "创建异常")
		}
		err = this.model.Add(channels)
		go this.createQrcode(channels)
		return err
}

func (this *popularizationServiceImpl) Incr(channel string, key string, incr ...int) error {
		var data = this.GetChannelInfo(channel)
		if data == nil {
				return common.NewErrors(common.NotFound, "渠道异常")
		}
		if len(incr) == 0 {
				incr = append(incr, 1)
		}
		var (
				numValid   = data.ValidNumber
				numInvited = data.InvitedRegisterNumber
		)
		if !data.Incr(key, incr...) {
				return common.NewErrors(common.ErrorCode, "incr异常")
		}
		var (
				where  = bson.M{"_id": data.Id, "status": data.Status, "invitedRegisterNumber": numInvited, "validNumber": numValid}
				update = bson.M{"invitedRegisterNumber": data.InvitedRegisterNumber, "validNumber": data.ValidNumber, "updatedAt": time.Now().Local()}
		)
		return this.model.Update(where, update)
}
