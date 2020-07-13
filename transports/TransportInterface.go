package transports

import (
		"fmt"
		"github.com/astaxie/beego"
		"reflect"
)

type TransportInterface interface {
		Empty() bool
		Init() TransportInterface
		M(filters ...func(m beego.M) beego.M) beego.M
}

type TransportBootstrap interface {
		Boot()
}

type TransportImpl struct {
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
func (this *TransportImpl) Init() TransportInterface {
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
func (this *TransportImpl) getBootstrap() TransportBootstrap {
		return this
}

// 相关函数执行次数
func (this *TransportImpl) Times(name string) int {
		if this._called == nil {
				return 0
		}
		return this._called[name]
}

// 增加
func (this *TransportImpl) IncrTimes(name string) TransportInterface {
		if this._called == nil {
				this._called = make(map[string]int)
		}
		this._called[name] = +1
		return this
}

// 引导加载
func (this *TransportImpl) Boot() {
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
func (this *TransportImpl) Empty() bool {
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
func (this *TransportImpl) AppendInit(handler func()) TransportInterface {
		this._init = append(this._init, handler)
		return this
}

//  注册的相关函数
func (this *TransportImpl) Register(fnName string, v interface{}) TransportInterface {
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
func (this *TransportImpl) GetHandler(key string) interface{} {
		return this._funcMap[key]
}

// 过滤输出
func (this *TransportImpl) M(filters ...func(m beego.M) beego.M) beego.M {
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
func (this *TransportImpl) SetPayLoad(m beego.M) TransportInterface {
		this._cache = &m
		return this
}

// 初始化数据解析
func (this *TransportImpl) InitPayload() *TransportImpl {
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
func (this *TransportImpl) GetPayLoad() beego.M {
		return *this._cache
}

// 添加过滤器｜格式器
func (this *TransportImpl) AddFilter(filters ...func(m beego.M) beego.M) TransportInterface {
		this._filters = append(this._filters, filters...)
		return this
}

// 销毁
func (this *TransportImpl) Destroy() {
		this._cache = nil
		this._init = nil
		this._funcMap = nil
		this._called = nil
		this._filters = nil
}

// 打印
func (this *TransportImpl) Dump() {
		if this.Empty() {
				fmt.Println("{}")
				return
		}
		fmt.Printf("%#v\n", this.GetPayLoad())
}

// 自定义打印
func (this *TransportImpl) Print() {
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
