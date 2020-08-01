package transports

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/transforms"
)

// 登陆参数
type Login struct {
		Mobile        string `json:"mobile"`
		Code          string `json:"code"`
		Username      string `json:"username"`
		Password      string `json:"password"`
		Email         string `json:"email"`
		transportImpl `json:",omitempty"`
}

func NewLoginInstance() *Login {
		var login = new(Login)
		return login
}

func NewLogin(ctx ...*context.BeegoInput) *Login {
		var login = NewLoginInstance()
		if len(ctx) > 0 {
				login.Load(ctx[0]).Init()
		}
		return login
}

func (this *Login) Boot() {
		this.Register("getPayload", this.getPayLoad)
		this.AddFilter(transforms.FilterEmpty)
}

func (this *Login) getPayLoad() error {
		this.SetPayLoad(beego.M{
				"mobile":   this.Mobile,
				"code":     this.Code,
				"username": this.Username,
				"password": this.Password,
				"email":    this.Email,
		})
		return nil
}

func (this *Login) Load(ctx *context.BeegoInput) *Login {
		var (
				_      = json.Unmarshal(ctx.RequestBody, this)
				mapper = map[string]interface{}{
						"mobile":   &this.Mobile,
						"code":     &this.Code,
						"username": &this.Username,
						"email":    &this.Email,
						"password": &this.Password,
				}
		)
		if this.Email == "" && this.Username == "" && this.Mobile == "" {
				for key, addr := range mapper {
						_ = ctx.Bind(addr, key)
				}
		}
		this.Init()
		return this
}

func (this *Login) Init() TransportInterface {
		this.Boot()
		this.transportImpl.Init()
		return this
}
