package repositories

import (
		"github.com/astaxie/beego/logs"
		"sync"
)

type EventRepository interface {
		Dispatch(target string, data interface{}, queue string)         // 派发事件到队列
		AddListener(target string, handler func(data interface{}) bool) // 添加事件监听器
		Once(target string, handler func(data interface{}) bool)        // 一次性事件
		Subscriber(subscriber Subscriber)                               //  注册订阅器
		Emit(target string, data interface{})                           // 发送事件
		Dispatcher(name string, dispatcher Dispatcher)                  // 添加事件派发器
}

type Subscriber interface {
		Target() string
		Register(repository EventRepository)
		Handle(data interface{}) bool
}

var (
		locker        sync.Once
		eventRegister *eventRepositoryImpl
)

func init() {
		if eventRegister == nil {
				locker.Do(newEventBus)
		}
}

func newEventBus() {
		eventRegister = new(eventRepositoryImpl)
		eventRegister.Init()
}

type eventRepositoryImpl struct {
		storage     map[string][]func(data interface{}) bool
		tmp         map[string][]func(data interface{}) bool
		dispatchers map[string][]Dispatcher
		locker      *sync.Mutex
}

type Dispatcher interface {
		Send(target string, data interface{})
}

func GetEventProvider() EventRepository {
		return eventRegister
}

func (this *eventRepositoryImpl) Init() {
		if this.storage == nil {
				this.storage = make(map[string][]func(data interface{}) bool)
		}
		if this.tmp == nil {
				this.tmp = make(map[string][]func(data interface{}) bool)
		}
		if this.locker == nil {
				this.locker = &sync.Mutex{}
		}
		if this.dispatchers == nil {
				this.dispatchers = map[string][]Dispatcher{}
		}
}

func (this *eventRepositoryImpl) Dispatch(target string, data interface{}, queue string) {
		dispatchers := this.getDispatcher(queue)
		if dispatchers == nil {
				// 缺失日志
				logs.Debug(target, data, queue)
				return
		}
		for _, dispatcher := range dispatchers {
				dispatcher.Send(target, data)
		}
}

func (this *eventRepositoryImpl) AddListener(target string, handler func(data interface{}) bool) {
		this.add(target, handler)
}

func (this *eventRepositoryImpl) Once(target string, handler func(data interface{}) bool) {
		this.locker.Lock()
		defer this.locker.Unlock()
		if _, ok := this.tmp[target]; !ok {
				this.tmp[target] = []func(interface{}) bool{}
		}
		this.tmp[target] = append(this.tmp[target], handler)
		return
}

func (this *eventRepositoryImpl) Subscriber(subscriber Subscriber) {
		if subscriber == nil {
				return
		}
		this.AddListener(subscriber.Target(), subscriber.Handle)
}

func (this *eventRepositoryImpl) Emit(target string, data interface{}) {
		this.locker.Lock()
		defer this.locker.Unlock()
		if handlers, ok := this.storage[target]; ok {
				for _, fn := range handlers {
						if !fn(data) {
								break
						}
				}
				return
		}
		if handlers, ok := this.tmp[target]; ok {
				for _, fn := range handlers {
						if !fn(data) {
								break
						}
				}
				delete(this.tmp, target)
				return
		}
}

func (this *eventRepositoryImpl) Dispatcher(name string, dispatcher Dispatcher) {
		if _, ok := this.dispatchers[name]; !ok {
				this.dispatchers[name] = []Dispatcher{}
		}
		this.locker.Lock()
		this.dispatchers[name] = append(this.dispatchers[name], dispatcher)
		this.locker.Unlock()
}

func (this *eventRepositoryImpl) getDispatcher(name string) []Dispatcher {
		if dispatchers, ok := this.dispatchers[name]; ok {
				return dispatchers
		}
		return nil
}

func (this *eventRepositoryImpl) add(target string, handler func(data interface{}) bool) *eventRepositoryImpl {
		this.locker.Lock()
		defer this.locker.Unlock()
		if _, ok := this.storage[target]; !ok {
				this.storage[target] = []func(interface{}) bool{}
		}
		this.storage[target] = append(this.storage[target], handler)
		return this
}
