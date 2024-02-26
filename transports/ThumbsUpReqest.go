package transports

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/transforms"
)

// ThumbsUpRequest 点赞请求
type ThumbsUpRequest struct {
		Type          string `json:"type"`
		Id            string `json:"id"`
		transportImpl `json:",omitempty"`
}

func (this *ThumbsUpRequest) Boot() {
		this.AddFilter(transforms.FilterEmpty)
		this.Register(getPayloadFn, this.getPayLoad)
}

func (this *ThumbsUpRequest) Load(data []byte) error {
		this.Init()
		return json.Unmarshal(data, this)
}

func NewThumbUpRequest() * ThumbsUpRequest {
		return new(ThumbsUpRequest)
}

func (this *ThumbsUpRequest) ParseFrom(ctx *context.BeegoInput) error {
		var (
				err    = json.Unmarshal(ctx.RequestBody, this)
				mapper = map[string]interface{}{
						"type": &this.Type,
						"id":   &this.Id,
				}
		)
		if err != nil {
				for key, addr := range mapper {
						err = ctx.Bind(addr, key)
				}
		}
		this.Init()
		return err
}

func (this *ThumbsUpRequest) getPayLoad() error {
		this.SetPayLoad(beego.M{
				"type": this.Type,
				"id":   this.Id,
		})
		return nil
}

func (this *ThumbsUpRequest)IsEmpty() bool  {
		if this.Id == "" {
				return true
		}
		return false
}

func (this *ThumbsUpRequest) GetType() string {
		if this.Type == "" {
				return "post"
		}
		return this.Type
}

func (this *ThumbsUpRequest) GetId() string {
		return this.Id
}

func (this *ThumbsUpRequest) GetObjectId() bson.ObjectId {
		if this.Id != "" {
				return bson.ObjectIdHex(this.Id)
		}
		return ""
}