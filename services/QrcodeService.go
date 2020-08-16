package services

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/siddontang/go/bson"
		"github.com/tuotoo/qrcode"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/plugins"
		"io"
		"os"
		"path/filepath"
		"strings"
)

type QrcodeService interface {
		CreateQrcode(data *QrcodeParams) (string, error)
}

type QrcodeParams struct {
		Content string               `json:"content"`
		Level   qrcode.RecoveryLevel `json:"level"`
		Options plugins.Options      `json:"options"`
		Params  beego.M              `json:"params"`
		UserId  string               `json:"userId"`
}

func NewQrcodeParams() *QrcodeParams {
		var param = new(QrcodeParams)
		param.Options.Init()
		if param.Level == 0 {
				param.Level = qrcode.RecoveryLevel(param.Options.Level)
		}
		return param
}

type qrcodeServiceImpl struct {
		BaseService
}

const (
		QrcodeParamSign    = "c=Q"
		QrcodeServiceClass = "QrcodeService"
)

func QrcodeServiceOf() QrcodeService {
		var service = new(qrcodeServiceImpl)
		service.Init()
		return service
}

func (this *qrcodeServiceImpl) Init() {
		this.ClassName = QrcodeServiceClass
		this.Constructor = func(args ...interface{}) interface{} {
				return QrcodeServiceOf()
		}
		this.init()
}

func (this *qrcodeServiceImpl) CreateQrcode(data *QrcodeParams) (string, error) {
		data.Set("mediaId", bson.NewObjectId().Hex())
		var (
				content    = data.getContent()
				provider   = this.GetProvider()
				userId     = data.Pop("userId")
				referId    = data.Pop("referId")
				referTable = data.Pop("referName", models.PopularizationChannelsTable)
				err        = provider.Save(data.addSignature(content), &data.Options)
		)
		if userId == "" {
				userId = data.UserId
		}
		var extras = beego.M{
				"userId":    userId,
				"referId":   referId,
				"referName": referTable,
				"_id":       bson.ObjectIdHex(data.GetParam("mediaId")),
		}
		return this.SaveQrcodeToUrl(data.Options.FileName, extras), err
}

func (this *qrcodeServiceImpl) SaveQrcodeToUrl(filename string, data ...beego.M) string {
		var (
				extras  = beego.M{}
				fs, err = os.Open(filename)
		)
		if err != nil {
				return ""
		}
		defer Close(fs)
		extras = beego.M{
				"fileType": AttachTypeImage,
				"filename": filepath.Base(fs.Name()),
		}
		extras = libs.MapMerge(extras, data[0])
		var attach = this.getAttachment().Save(fs, extras)
		if attach == nil {
				return ""
		}
		return attach.GetUrl()
}

func (this *qrcodeServiceImpl) getAttachment() AttachmentService {
		return AttachmentServiceOf()
}

func (this *qrcodeServiceImpl) GetProvider() *plugins.Qrcode {
		return plugins.GetQrcode()
}

func (this *QrcodeParams) Init() {
		this.Content = ""
		this.Params = beego.M{}
}

func (this *QrcodeParams) getContent() string {
		if len(this.Params) == 0 {
				if this.Content == "" {
						return ""
				}
				return this.Content
		}
		var content = this.Content
		for k, v := range this.Params {
				value := k + "=" + fmt.Sprintf("%v", v)
				if content == "" {
						content = value
						continue
				}
				content = "&" + value
		}
		return content
}

func (this *QrcodeParams) Verify(sign string) bool {
		var data = this.ParseContent(this.getContent())
		delete(data, "s")
		this.Params = data
		this.Content = ""
		return this.signature(this.getContent()) == sign
}

func (this *QrcodeParams) signature(content string) string {
		return libs.Encrypt(content + "&" + QrcodeParamSign)
}

func (this *QrcodeParams) addSignature(content string) string {
		return content + "&s=" + libs.Encrypt(content+"&"+QrcodeParamSign)
}

func (this *QrcodeParams) ParseContent(content string) beego.M {
		if content == "" {
				return beego.M{}
		}
		var (
				result = beego.M{}
		)
		for _, v := range strings.SplitN(content, "&", -1) {
				if !strings.Contains(v, "=") {
						continue
				}
				arr := strings.SplitN(v, "=", 2)
				if len(arr) >= 2 {
						result[strings.TrimSpace(arr[0])] = strings.TrimSpace(arr[1])
				}
		}
		return result
}

func (this *QrcodeParams) Pop(key string, defaults ...string) string {
		if len(defaults) == 0 {
				defaults = append(defaults, "")
		}
		var v, ok = this.Params[key]
		if !ok {
				return defaults[0]
		}
		delete(this.Params, key)
		return v.(string)
}

func (this *QrcodeParams) Set(key string, v string) *QrcodeParams {
		this.Params[key] = v
		return this
}

func (this *QrcodeParams) GetParam(key string, defaults ...string) string {
		if len(defaults) == 0 {
				defaults = append(defaults, "")
		}
		var v, ok = this.Params[key]
		if !ok {
				return defaults[0]
		}
		return v.(string)
}

func Close(closer io.Closer) {
		var err = closer.Close()
		if err != nil {
				logs.Error(err)
		}
}
