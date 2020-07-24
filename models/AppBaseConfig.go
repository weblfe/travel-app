package models

import (
		"errors"
		"github.com/astaxie/beego"
)

// App 配置信息
type AppBaseConfig struct {
		AppCode          string `json:"appCode"`          // 授权码
		Download         string `json:"download"`         // 下载页面
		AppHelpUrl       string `json:"appHelpUrl"`       // 帮助页
		AppPrivacy       string `json:"appPrivacy"`       // 隐私条款页
		AppAboutUs       string `json:"appAboutUs"`       // 关于我们
		AppUserAgreement string `json:"appUserAgreement"` // 用户协议页
		Email            string `json:"email"`            // 客服邮箱
		QQ               string `json:"QQ"`               // 客服QQ
		WeChat           string `json:"wechat"`           // 微信
		dataClassImpl    `bson:",omitempty"  json:",omitempty"`
}

func (this *AppBaseConfig) Init() {
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

// app 基础配置数据
func (this *AppBaseConfig) data() beego.M {
		return beego.M{
				"appCode":          this.GetAppCode(),
				"download":         this.GetDownloadUrl(),
				"appHelpUrl":       this.GetHelpUrl(),
				"appPrivacy":       this.GetAppPrivacy(),
				"appAboutUs":       this.GetAppAboutUs(),
				"appUserAgreement": this.GetAppUserAgreement(),
				"email":            this.GetMail(),
				"QQ":               this.GetQQ(),
				"wechat":           this.GetWeChat(),
		}
}

func (this *AppBaseConfig) GetDownloadUrl(drivers ...string) string {
		if this.Download != "" {
				return this.Download
		}
		if len(drivers) == 0 {
				drivers = append(drivers, "")
		}
		return AppModelOf().GetAppDownloadUrl(drivers[0])
}

func (this *AppBaseConfig) GetHelpUrl(drivers ...string) string {
		if this.Download != "" {
				return this.Download
		}
		if len(drivers) == 0 {
				drivers = append(drivers, DefaultConfigRoot)
		}
		return this.get("appHelpUrl", drivers...)
}

func (this *AppBaseConfig) GetAppAboutUs(drivers ...string) string {
		if this.Download != "" {
				return this.Download
		}
		if len(drivers) == 0 {
				drivers = append(drivers, DefaultConfigRoot)
		}
		return this.get("appAboutUs", drivers...)
}

func (this *AppBaseConfig) GetAppUserAgreement(drivers ...string) string {
		if this.Download != "" {
				return this.Download
		}
		if len(drivers) == 0 {
				drivers = append(drivers, DefaultConfigRoot)
		}
		return this.get("appUserAgreement", drivers...)
}

func (this *AppBaseConfig) GetMail(drivers ...string) string {
		if this.Download != "" {
				return this.Download
		}
		if len(drivers) == 0 {
				drivers = append(drivers, DefaultConfigRoot)
		}
		return this.get("email", drivers...)
}

func (this *AppBaseConfig) GetWeChat(drivers ...string) string {
		if this.Download != "" {
				return this.Download
		}
		if len(drivers) == 0 {
				drivers = append(drivers, DefaultConfigRoot)
		}
		return this.get("wechat", drivers...)
}

func (this *AppBaseConfig) GetQQ(drivers ...string) string {
		if this.Download != "" {
				return this.Download
		}
		if len(drivers) == 0 {
				drivers = append(drivers, DefaultConfigRoot)
		}
		return this.get("QQ", drivers...)
}

func (this *AppBaseConfig) GetAppPrivacy(drivers ...string) string {
		if this.Download != "" {
				return this.Download
		}
		if len(drivers) == 0 {
				drivers = append(drivers, DefaultConfigRoot)
		}
		return this.get("appPrivacy", drivers...)
}

// 获取AppCode
func (this *AppBaseConfig) GetAppCode(drivers ...string) string {
		if this.AppCode != "" {
				return this.AppCode
		}
		if len(drivers) == 0 {
				drivers = append(drivers, DefaultConfigRoot)
		}
		return this.get("appCode", drivers...)
}

func (this *AppBaseConfig) get(key string, drivers ...string) string {
		var value = ConfigModelOf().GetString(key, drivers[0])
		if value != "" {
				return value
		}
		// app 未指定设备类型
		if drivers[0] == DefaultConfigRoot {
				return ""
		}
		var name = ConfigModelOf().GetString(key)
		if name == "" {
				return DefaultAppName
		}
		return name
}

func (this *AppBaseConfig) Load(driver ...string) *AppBaseConfig {
		this.Init()
		this.init(driver...)
		return this
}

func (this *AppBaseConfig) init(driver ...string) {

		if this.AppCode == "" {
				this.AppCode = this.GetAppCode(driver...)
		}
		if this.WeChat == "" {
				this.WeChat = this.GetWeChat(driver...)
		}
		if this.QQ == "" {
				this.QQ = this.GetQQ(driver...)
		}
		if this.Email == "" {
				this.Email = this.GetMail(driver...)
		}
		if this.AppUserAgreement == "" {
				this.AppUserAgreement = this.GetAppUserAgreement(driver...)
		}
		if this.Download == "" {
				this.Download = this.GetDownloadUrl(driver...)
		}
		if this.AppAboutUs == "" {
				this.AppAboutUs = this.GetAppAboutUs(driver...)
		}
		if this.AppHelpUrl == "" {
				this.AppHelpUrl = this.GetHelpUrl(driver...)
		}
		if this.AppPrivacy == "" {
				this.AppPrivacy = this.GetAppPrivacy(driver...)
		}
}

func (this *AppBaseConfig) setDefaults() {
		if this.AppCode == "" {
				this.AppCode = this.GetAppCode()
		}
		if this.WeChat == "" {
				this.WeChat = this.GetWeChat()
		}
		if this.QQ == "" {
				this.QQ = this.GetQQ()
		}
		if this.Email == "" {
				this.Email = this.GetMail()
		}
		if this.AppUserAgreement == "" {
				this.AppUserAgreement = this.GetAppUserAgreement()
		}
		if this.AppAboutUs == "" {
				this.AppAboutUs = this.GetAppAboutUs()
		}
		if this.AppPrivacy == "" {
				this.AppPrivacy = this.GetAppPrivacy()
		}

}

func (this *AppBaseConfig) setAttributes(data map[string]interface{}, safe ...bool) {
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

func (this *AppBaseConfig) Set(key string, v interface{}) *AppBaseConfig {
		switch key {
		case "appCode":
				this.SetString(&this.AppCode, v)
		case "download":
				this.SetString(&this.Download, v)
		case "appHelpUrl":
				this.SetString(&this.AppHelpUrl, v)
		case "appPrivacy":
				this.SetString(&this.AppPrivacy, v)
		case "appAboutUs":
				this.SetString(&this.AppAboutUs, v)
		case "appUserAgreement":
				this.SetString(&this.AppUserAgreement, v)
		case "email":
				this.SetString(&this.Email, v)
		case "QQ":
				this.SetString(&this.QQ, v)
		case "wechat":
				this.SetString(&this.WeChat, v)
		}
		return this
}

// @todo
func (this *AppBaseConfig) save() error {
		return errors.New("不支持批量更新")
}
