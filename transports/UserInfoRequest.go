package transports

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/transforms"
)

// 请求参数
type UpdateUserRequest struct {
		AvatarId          string   `json:"avatarId,omitempty"`
		NickName          string   `json:"nickname,omitempty"`
		Email             string   `json:"email,omitempty"`
		Gender            int      `json:"gender,omitempty"`
		Intro             string   `json:"intro,omitempty"`
		BackgroundCoverId string   `json:"backgroundCoverId,omitempty"`
		Modifies          []string `json:"modifies,omitempty"`
		TransportImpl
}

func (this *UpdateUserRequest) Boot() {
		this.Register("getPayload", this.getPayLoad)
		this.AddFilter(this.filter, transforms.FilterEmpty)
}

func (this *UpdateUserRequest) Load(data []byte) error {
		this.Init()
		return json.Unmarshal(data, this)
}

func (this *UpdateUserRequest) getPayLoad() error {
		this.SetPayLoad(beego.M{
				"avatarId":          this.AvatarId,
				"nickname":          this.NickName,
				"email":             this.Email,
				"gender":            this.Gender,
				"intro":             this.Intro,
				"modifies":          this.Modifies,
				"backgroundCoverId": this.BackgroundCoverId,
		})
		return nil
}

func (this *UpdateUserRequest) filter(data beego.M) beego.M {
		if this.Gender == 0 {
				delete(data, "gender")
		}
		if len(this.Modifies) == 0 {
				delete(data, "modifies")
		}
		return data
}

func (this *UpdateUserRequest) Init() TransportInterface {
		this.Boot()
		this.TransportImpl.Init()
		return this
}

// 重置密码请求体
type ResetPassword struct {
		Password        string `json:"password"`                  // 新密码
		CurrentPassword string `json:"currentPassword,omitempty"` // 当前登陆使用的密码
		UserId          string `json:"userId,omitempty"`          // 当前用户ID
		Code            string `json:"code,omitempty"`            // 手机重置密码使用的验证码
		Mobile          string `json:"mobile,omitempty"`          // 手机号
		TransportImpl
}

func (this *ResetPassword) Load(ctx *context.BeegoInput) *ResetPassword {
		var (
				_      = json.Unmarshal(ctx.RequestBody, this)
				mapper = map[string]interface{}{
						"code":            &this.Code,
						"mobile":          &this.Mobile,
						"password":        &this.Password,
						"currentPassword": &this.CurrentPassword,
				}
		)
		if this.Password == "" {
				for key, addr := range mapper {
						_ = ctx.Bind(addr, key)
				}
		}
		this.Init()
		return this
}

func (this *ResetPassword) getPayLoad() error {
		this.SetPayLoad(beego.M{
				"userId":          this.UserId,
				"password":        this.Password,
				"currentPassword": this.CurrentPassword,
				"code":            this.Code,
				"mobile":          this.Mobile,
		})
		return nil
}

func (this *ResetPassword) Boot() {
		this.Register("getPayload", this.getPayLoad)
		this.AddFilter(transforms.FilterEmpty)
}

func (this *ResetPassword) Init() TransportInterface {
		this.Boot()
		this.TransportImpl.Init()
		return this
}
