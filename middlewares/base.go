package middlewares

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/weblfe/travel-app/libs"
		"sync"
)

type Middleware interface {
		Register(args ...interface{})
		GetHandler() beego.FilterFunc
		Middleware() string
}

type middlewareImpl struct {
		handler beego.FilterFunc
		Pattern string
		Pos     int
		Name    string
}

type Initialization interface {
		Init()
}

func (this *middlewareImpl) GetHandler() beego.FilterFunc {
		return this.handler
}

func (this *middlewareImpl) SetHandler(handlers ...func(ctx *context.Context) bool) Middleware {
		var argc = len(handlers)
		if this.handler == nil && argc > 0 {
				var (
						i       = 0
						fn      beego.FilterFunc
						handler = handlers[i]
				)
				if argc == 1 {
						handlers = append(handlers, nil)
				}
				for {
						if i+1 == argc {
								fn = this.Wrapper(handler, handlers[i+1])
								break
						}
						if i+1 < argc {
								handler = this.FilterWrapper(handler, handlers[i+1])
								i++
						}
						i++
						if i > argc {
								break
						}
				}
				this.handler = fn
		}
		return this
}

func (this *middlewareImpl) FilterWrapper(fn func(ctx *context.Context) bool, next func(ctx *context.Context) bool) func(ctx *context.Context) bool {
		return func(ctx *context.Context) bool {
				if !fn(ctx) {
						return false
				}
				if next == nil {
						return true
				}
				return next(ctx)
		}
}

func (this *middlewareImpl) Wrapper(fn func(ctx *context.Context) bool, next func(ctx *context.Context) bool) func(ctx *context.Context) {
		return func(ctx *context.Context) {
				if !fn(ctx) {
						return
				}
				if next == nil {
						return
				}
				next(ctx)
		}
}

func (this *middlewareImpl) Register(args ...interface{}) {
		var (
				argc = len(args)
		)
		if argc == 0 {
				args[0] = this.Pattern
		}
		if argc < 2 {
				args[1] = this.Pos
		}
		beego.InsertFilter(args[0].(string), args[1].(int), this.GetHandler())
}

func (this *middlewareImpl) Middleware() string {
		return this.Name
}

type middlewareManager struct {
		container map[string]Middleware
		mutex     *sync.Mutex
		routers   map[string]*middlewareEntry
}

type middlewareEntry struct {
		Pattern string
		Pos     int
		Mid     Middleware
}

var (
		locks          map[string]*sync.Once
		instanceManger *middlewareManager
)

const (
		MiddlewareManger = "manager"
)

func init() {
		if instanceManger == nil {
				getLock(MiddlewareManger).Do(func() {
						instanceManger = new(middlewareManager)
						instanceManger.Init()
				})
		}
}

func getLock(name string) *sync.Once {
		if locks == nil {
				locks = make(map[string]*sync.Once)
		}
		if l, ok := locks[name]; ok {
				return l
		}
		locks[name] = &sync.Once{}
		return locks[name]
}

func GetMiddleware(name string) Middleware {
		return instanceManger.Get(name)
}

func Register(name string, mid Middleware) *middlewareManager {
		return instanceManger.Set(name, mid)
}

func GetMiddlewareManger() *middlewareManager {
		return instanceManger
}

func (this *middlewareManager) Init() {
		if this.container == nil {
				this.container = make(map[string]Middleware)
		}
		if this.routers == nil {
				this.routers = make(map[string]*middlewareEntry)
		}
		if this.mutex == nil {
				this.mutex = &sync.Mutex{}
		}
}

func (this *middlewareManager) Get(name string) Middleware {
		return this.container[name]
}

func (this *middlewareManager) Set(name string, mid Middleware) *middlewareManager {
		this.mutex.Lock()
		this.container[name] = mid
		this.mutex.Unlock()
		return this
}

func (this *middlewareManager) Exists(name string) bool {
		if _, ok := this.container[name]; ok {
				return true
		}
		return false
}

func (this *middlewareManager) Router(mid, pathPattern string, pos int) *middlewareManager {
		var (
				entry = new(middlewareEntry)
				key   = libs.HashCode(mid + pathPattern + fmt.Sprintf("%d", pos))
		)
		if !this.Exists(mid) {
				return this
		}
		if _, ok := this.routers[key]; ok {
				return this
		}
		entry.Pos = pos
		entry.Pattern = pathPattern
		entry.Mid = this.Get(mid)
		this.routers[key] = entry
		return this
}

func (this *middlewareManager) Boot() {
		for key, it := range this.routers {
				it.Mid.Register(it.Pattern, it.Pos)
				delete(this.routers, key)
		}
}
