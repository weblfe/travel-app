package transports

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
)

// 验证码参数
type MobileRequest struct {
		Mobile string `json:"mobile"`
		Type   string `json:"type"`
		transportImpl
}



// 加载
func (this *MobileRequest) Load(ctx *context.BeegoInput) *MobileRequest {
		if err := ctx.Bind(&this.Mobile, "mobile"); err != nil || this.Mobile == "" {
				_ = json.Unmarshal(ctx.RequestBody, this)
		}
		if this.Type == "" {
				_ = ctx.Bind(&this.Type, "type")
		}
		this.Init()
		return this
}

func (this *MobileRequest) getPayLoad() error {
		this.SetPayLoad(beego.M{
				"mobile": this.Mobile,
				"type":   this.Type,
		})
		return nil
}

func (this *MobileRequest)Init() TransportInterface {
		this.Boot()
		this.transportImpl.Init()
		return this
}

func (this *MobileRequest) Boot() {
		this.Register(GetPayload, this.getPayLoad)
}
