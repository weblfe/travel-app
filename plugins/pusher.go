package plugins

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/nats-io/nats.go"
		"reflect"
		"time"
)

type PusherInterface interface {
		GetChannel() string
		Push(string, interface{}) error
		SetChannel(string) PusherInterface
		Subscribe(string, interface{}) (*nats.Subscription, error)
		Close()
}

type MsgRpc interface {
		Close()
		GetTimeOut() time.Duration
		SetTimeOut(duration time.Duration) MsgRpc
		Worker(handler interface{}) (*nats.Subscription, error)
		RpcCall(method string, data interface{}, responsePtr interface{}, timeOut ...time.Duration) error
		RpcService(method string, handler func(string, []byte) (interface{}, error)) (*nats.Subscription, error)
}

type MapperDecoder interface {
		M(...func(beego.M)) beego.M
}

type Pusher struct {
		Client  *NatsPlugin
		channel string
		conn    *nats.EncodedConn
		Timeout time.Duration
}

func GetPusher() PusherInterface {
		var instance = new(Pusher)
		instance.Init()
		return instance
}

func GetMsgRpc() MsgRpc {
		var instance = new(Pusher)
		instance.Init()
		return instance
}

func (this *Pusher) Init() {
		this.Client = GetNatsPlugin()
}

func (this *Pusher) Push(channel string, data interface{}) error {
		return this.client().Publish(channel, data)
}

func (this *Pusher) SetChannel(channel string) PusherInterface {
		if this.channel == "" {
				this.channel = channel
		}
		return this
}

func (this *Pusher) Send(data []byte) error {
		return this.Push(this.GetChannel(), data)
}

func (this *Pusher) GetChannel() string {
		if this.channel == "" {
				return "default"
		}
		return this.channel
}

// 订阅
func (this *Pusher) Subscribe(channel string, handler interface{}) (*nats.Subscription, error) {
		return this.client().Subscribe(channel, handler)
}

// 队列订阅
func (this *Pusher) QueueSubscribe(channel, queue string, cb interface{}) (*nats.Subscription, error) {
		return this.client().QueueSubscribe(channel, queue, cb)
}

// 队列订阅
func (this *Pusher) Queue(channel, queue string, cb interface{}) (*nats.Subscription, error) {
		return this.client().BindRecvQueueChan(channel, queue, cb)
}

// worker
func (this *Pusher) Worker(handler interface{}) (*nats.Subscription, error) {
		return this.Subscribe(this.GetChannel(), handler)
}

// 调用远程服务
func (this *Pusher) RpcService(method string, handler func(string, []byte) (interface{}, error)) (*nats.Subscription, error) {
		return this.Subscribe(method, func(msg *nats.Msg) error {
				var data, err = handler(msg.Subject, msg.Data)
				if err == nil {
						return respond(msg, data)
				}
				return respond(msg, err.Error())
		})
}

// 远程调用
func (this *Pusher) RpcCall(method string, data interface{}, responsePtr interface{}, timeOut ...time.Duration) error {
		if len(timeOut) == 0 {
				timeOut = append(timeOut, this.GetTimeOut())
		}
		return this.client().Request(method, data, responsePtr, timeOut[0])
}

func (this *Pusher) Close() {
		this.client().Close()
		this.Client = nil
}

func (this *Pusher) GetTimeOut() time.Duration {
		if this.Timeout != 0 {
				return this.Timeout
		}
		return 2 * time.Second
}

func (this *Pusher) SetTimeOut(duration time.Duration) MsgRpc {
		if duration <= 0 {
				return this
		}
		this.Timeout = duration
		return this
}

func (this *Pusher) client() *nats.EncodedConn {
		if this.conn != nil {
				return this.conn
		}
		this.conn, _ = nats.NewEncodedConn(this.Client.GetConn(), nats.JSON_ENCODER)
		return this.conn
}

func respond(msg *nats.Msg, v interface{}) error {
		switch v.(type) {
		case string:
				return msg.Respond([]byte(v.(string)))
		case []byte:
				return msg.Respond(v.([]byte))
		}
		var (
				getType  = reflect.TypeOf(v)
				getValue = reflect.ValueOf(v)
		)
		// func
		if getType.Kind() == reflect.Func {
				var (
						params = []reflect.Value{
								reflect.ValueOf(msg), reflect.ValueOf(v),
						}
						values = getValue.Call(params)
				)
				if e, ok := values[0].Interface().(error); ok {
						return e
				}
				return nil
		}
		// struct
		if getType.Kind() == reflect.Struct {
				var data, _ = json.Marshal(v)
				return msg.Respond(data)
		}
		// mapper
		if getType.Kind() == reflect.Map {
				var data, _ = json.Marshal(v)
				return msg.Respond(data)
		}
		if decoder, ok := v.(MapperDecoder); ok {
				var data, _ = json.Marshal(decoder.M())
				return msg.Respond(data)
		}
		return msg.Respond(NewErrorResponse())
}

func NewErrorResponse() []byte {
		var data, _ = json.Marshal(beego.M{
				"msg":  "service error",
				"code": 5000,
				"error": beego.M{
						"errno":  5000,
						"errmsg": "服务异常",
				},
				"data": nil,
		})
		return data
}
