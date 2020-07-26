package test

import (
		"encoding/json"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/plugins"
		"os"
		"testing"
)

func TestPusher(t *testing.T) {
		var (
				err    error
				text   = beego.M{
						"hello":"rpc",
				}
				data   = make([]byte, 100)
				pusher = plugins.GetMsgRpc()
		)
		data = data[:0]
		go func() {
				// worker
				_, _ = pusher.RpcService("hello", func(s string, bytes []byte) (interface{}, error) {
						var data = beego.M{}
						_ = json.Unmarshal(bytes,&data)
						data["method"] = s
						return data, nil
				})
		}()
		// push
		err = pusher.RpcCall("hello", text, &data)
		if err == nil {
				fmt.Println("call:", string(data))
		} else {
				fmt.Println("err:", err.Error())
		}

}

func init() {
		var err error
		err = os.Setenv("NATS_CONN_URL", "127.0.0.1:4222")
		err = os.Setenv("NATS_NKEY", "UAY5N6KNU7EMGIHIMG4XXMBZDLXGYVPVEPP5KL2DAH3JXMYNGGNPCUT4")
		err = os.Setenv("NATS_SEED_KEY", "SUACWCQ2CU5HOBVM3IXMV4JO6HG2L5KTGCXQSR3AG7CRYQXWSA7ECWINMY")
		if err == nil {
				return
		}
		fmt.Println(err.Error())
}
