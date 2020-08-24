package transports

import (
		"bytes"
		"encoding/gob"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/libs"
		"net/http"
		"net/url"
		"reflect"
		"strings"
)

type TransportInterface interface {
		Empty() bool
		Init() TransportInterface
		M(filters ...func(m beego.M) beego.M) beego.M
}

type TransportBootstrap interface {
		Boot()
}

type transportImpl struct {
		_cache   *beego.M
		_init    []func()
		_called  map[string]int
		_funcMap map[string]interface{}
		_filters []func(m beego.M) beego.M
}

const (
		GetPayload = "getPayload"
)

// 初始化
func (this *transportImpl) Init() TransportInterface {
		if this.Times("Init") > 0 {
				return this
		}
		if this._cache == nil {
				this._cache = new(beego.M)
		}
		if this._called == nil {
				this._called = make(map[string]int)
		}
		if len(this._funcMap) == 0 {
				this._funcMap = make(map[string]interface{})
		}
		if len(this._init) == 0 {
				this._init = make([]func(), 2)
				this._init = this._init[:0]
		}
		if len(this._filters) == 0 {
				this._filters = make([]func(m beego.M) beego.M, 2)
				this._filters = this._filters[:0]
		}
		this.getBootstrap().Boot()
		this._called["Init"] = +1
		return this
}

// 获取引导器
func (this *transportImpl) getBootstrap() TransportBootstrap {
		return this
}

// 相关函数执行次数
func (this *transportImpl) Times(name string) int {
		if this._called == nil {
				return 0
		}
		return this._called[name]
}

// 增加
func (this *transportImpl) IncrTimes(name string) TransportInterface {
		if this._called == nil {
				this._called = make(map[string]int)
		}
		this._called[name] = +1
		return this
}

// 引导加载
func (this *transportImpl) Boot() {
		if this.Times("Boot") > 0 {
				return
		}
		for _, init := range this._init {
				if init == nil {
						continue
				}
				init()
		}
		this._called["Boot"] = +1
}

// 数据是否为空
func (this *transportImpl) Empty() bool {
		var fn = this.GetHandler("empty")
		if fn == nil {
				if len(*this._cache) == 0 {
						this.M()
						return len(this.GetPayLoad()) == 0
				}
				return true
		}
		if handler, ok := fn.(func() bool); ok && handler != nil {
				return handler()
		}
		return false
}

// 添加初始相关函数
func (this *transportImpl) AppendInit(handler func()) TransportInterface {
		this._init = append(this._init, handler)
		return this
}

//  注册的相关函数
func (this *transportImpl) Register(fnName string, v interface{}) TransportInterface {
		if reflect.TypeOf(v).Kind() != reflect.Func {
				return this
		}
		if this._funcMap == nil {
				this._funcMap = make(map[string]interface{})
		}
		if _, ok := this._funcMap[fnName]; ok {
				return this
		}
		this._funcMap[fnName] = v
		return this
}

// 获取注册的相关函数
func (this *transportImpl) GetHandler(key string) interface{} {
		return this._funcMap[key]
}

// 过滤输出
func (this *transportImpl) M(filters ...func(m beego.M) beego.M) beego.M {
		var data = this.GetPayLoad()
		if len(data) == 0 {
				data = this.InitPayload().GetPayLoad()
		}
		// 过滤器|格式器
		if len(filters) == 0 {
				filters = this._filters
		} else {
				filters = append(this._filters, filters...)
		}
		for _, fn := range filters {
				if fn == nil {
						continue
				}
				data = fn(data)
		}
		return data
}

// 设置解析结果数据
func (this *transportImpl) SetPayLoad(m beego.M) TransportInterface {
		this._cache = &m
		return this
}

// 初始化数据解析
func (this *transportImpl) InitPayload() *transportImpl {
		var (
				n = this._called["getPayload"]
				v = this.GetHandler("getPayload")
		)
		if v == nil || n > 0 {
				return this
		}
		if fn, ok := v.(func() error); ok {
				err := fn()
				this._called["getPayload"] = +1
				if err == nil {
						return this
				}
		}
		return this
}

// 获取最原始解析处理的数据
func (this *transportImpl) GetPayLoad() beego.M {
		return *this._cache
}

// 添加过滤器｜格式器
func (this *transportImpl) AddFilter(filters ...func(m beego.M) beego.M) TransportInterface {
		this._filters = append(this._filters, filters...)
		return this
}

// 销毁
func (this *transportImpl) Destroy() {
		this._cache = nil
		this._init = nil
		this._funcMap = nil
		this._called = nil
		this._filters = nil
}

// 打印
func (this *transportImpl) Dump() {
		if this.Empty() {
				fmt.Println("{}")
				return
		}
		fmt.Printf("%#v\n", this.GetPayLoad())
}

// 自定义打印
func (this *transportImpl) Print() {
		if this.Empty() {
				fmt.Println("{}")
				return
		}
		fmt.Println("{")
		for key, v := range this.GetPayLoad() {
				fmt.Printf("  %s:%v ,\n", key, v)
		}
		fmt.Println("}")
}

// 数据对象复制
func (this *transportImpl) Copy(data interface{}, dest interface{}) error {
		var (
				err  error
				buff = new(bytes.Buffer)
				enc  = gob.NewEncoder(buff)
				dec  = gob.NewDecoder(buff)
		)
		err = enc.Encode(data)
		if err == nil {
				err = dec.Decode(dest)
				if err == nil {
						return nil
				}
		}
		return err
}

// 对象拷贝
func (this *transportImpl) Clone(source interface{}, dest interface{}) error {
		var byt, err = libs.Json().Marshal(source)
		if err == nil {
				err = libs.Json().Unmarshal(byt, dest)
				if err == nil {
						return nil
				}
		}
		return err
}

func (this *transportImpl) IsJson(header http.Header) bool {
		return strings.Contains(header.Get("Content-Type"), "json")
}

func (this *transportImpl) IsForm(header http.Header) bool {
		return strings.Contains(header.Get("Content-Type"), "form")
}

func (this *transportImpl) Decoder(ctx *context.BeegoInput, v interface{}) error {
		var (
				err  error
				data = beego.M{}
		)
		// 是否Get
		if this.IsGet(ctx.Context.Request.Method) {
				data := ctx.Context.Input.Params()
				return this.Copy(data, v)
		}
		// 是json 数据请求
		if this.IsJson(ctx.Context.Request.Header) {
				err = libs.Json().Unmarshal(ctx.RequestBody, v)
		}
		// 表单请求
		if err != nil || this.IsForm(ctx.Context.Request.Header) {
				if len(ctx.Context.Request.PostForm) > 0 {
						data = this.listFrom(ctx.Context.Request.PostForm)
				}
				if len(data) <= 0 && len(ctx.Context.Request.Form) > 0 {
						data = this.listFrom(ctx.Context.Request.Form)
				}
				return this.Clone(data, v)
		}
		return err
}

func (this *transportImpl) IsGet(method string) bool {
		return strings.EqualFold("GET", method)
}

func (this *transportImpl) listFrom(values url.Values) beego.M {
		var data = beego.M{}
		for key, items := range values {
				size := len(items)
				if size <= 0 {
						data[key] = nil
						continue
				}
				if size == 1 {
						data[key] = items[0]
						continue
				}
				data[key] = items
		}
		return data
}
