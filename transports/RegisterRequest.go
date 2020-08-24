package transports

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/transforms"
)

// 注册参数
type RegisterRequest struct {
		Code     string `json:"code"`
		Mobile   string `json:"mobile"`
		Account  string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Third    string `json:"third"`
		Type     string `json:"type"`
		Way      string `json:"_register"`
}

// 载入数据
func (this *RegisterRequest) Load(ctx *context.BeegoInput) *RegisterRequest {
		var (
				_      = libs.Json().Unmarshal(ctx.RequestBody, this)
				mapper = map[string]interface{}{
						"mobile":    &this.Mobile,
						"password":  &this.Password,
						"username":  &this.Account,
						"code":      &this.Code,
						"email":     &this.Email,
						"third":     &this.Third,
						"type":      &this.Type,
						"_register": &this.Way,
				}
		)
		for key, v := range mapper {
				if str, ok := v.(*string); ok {
						if *str != "" {
								continue
						}
				}
				_ = ctx.Bind(v, key)
		}

		return this
}

// 过滤
func (this *RegisterRequest) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = beego.M{
				"mobile":       this.Mobile,
				"passwordHash": this.Password,
				"username":     this.Account,
				"code":         this.Code,
				"email":        this.Email,
				"third":        this.Third,
				"type":         this.Type,
				"registerWay":  this.Way,
		}
		filters = append(filters, transforms.FilterEmpty)
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}

// 快捷注册
type QuickRegister struct {
		Mobile      string `json:"mobile"`
		Password    string `json:"password"`
		RegisterWay string `json:"_register"`
		Gender      int    `json:"gender"`
		UserName    string `json:"username"`
		NickName    string `json:"nickname"`
		Email       string `json:"email"`
		AvatarId    string `json:"avatarId"`
}

// 载入
func (this *QuickRegister) Load(ctx *context.BeegoInput) *QuickRegister {
		var (
				_      = libs.Json().Unmarshal(ctx.RequestBody, this)
				mapper = map[string]interface{}{
						"mobile":    &this.Mobile,
						"password":  &this.Password,
						"_register": &this.RegisterWay,
						"gender":    &this.Gender,
						"username":  &this.UserName,
						"nickname":  &this.NickName,
						"email":     &this.Email,
						"avatarId":  &this.AvatarId,
				}
		)
		if this.Mobile == "" {
				for key, v := range mapper {
						_ = ctx.Bind(v, key)
				}
		}
		return this
}

// 过滤
func (this *QuickRegister) M(filters ...func(m beego.M) beego.M) beego.M {
		if this.RegisterWay == "" {
				this.RegisterWay = "quick"
		}
		var data = beego.M{
				"mobile":       this.Mobile,
				"passwordHash": this.Password,
				"registerWay":  this.RegisterWay,
				"gender":       this.Gender,
				"username":     this.UserName,
				"nickname":     this.NickName,
				"email":        this.Email,
				"avatarId":     this.AvatarId,
		}
		filters = append(filters, transforms.FilterEmpty)
		for _, filter := range filters {
				data = filter(data)
		}
		return data
}

