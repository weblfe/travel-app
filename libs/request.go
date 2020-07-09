package libs

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
)

func LoadRequest(ctx *context.Context, request interface{}) error {
		if m, ok := request.(*beego.M); ok {
				return json.Unmarshal(ctx.Input.RequestBody, m)
		}

		if m, ok := request.(beego.M); ok && len(m) > 0 {
				for key, v := range m {
						if err := ctx.Input.Bind(v, key); err != nil {
								return err
						}
				}
		}
		return nil
}

func GetRequestMapper(ctx *context.Context) beego.M {
		var data = make(beego.M)
		if err := LoadRequest(ctx, &data); err == nil {
				return data
		}
		return data
}

