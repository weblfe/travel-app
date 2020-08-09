package libs

import (
		"errors"
		"sync"
)

type ContainerInterface interface {
		Exists(id string) bool
		Reset(name ...string) ContainerInterface
		Get(id string, args ...interface{}) interface{}
		Register(name string, factory func(args ...interface{}) interface{}, force ...bool) ContainerInterface
		GetInstance(name string, args ...interface{}) interface{}
}

type item struct {
		Index int
		Value interface{}
}

type container struct {
		locker    sync.Mutex
		cache     map[string]*item
		factories []func(args ...interface{}) interface{}
}

var (
		_container           *container
		_containerSyncLocker = sync.Once{}
		ContainerNotFound    = errors.New("not exists in container")
)

func Container() ContainerInterface {
		if _container == nil {
				_containerSyncLocker.Do(newContainer)
		}
		return _container
}

func newContainer() {
		_container = new(container)
		_container.init()
}

func (this *container) Lock() {
		this.locker.Lock()
}

func (this *container) UnLock() {
		this.locker.Unlock()
}

func (this *container) Exists(id string) bool {
		if _, ok := this.cache[id]; ok {
				return true
		}
		return false
}

func (this *container) init() ContainerInterface {
		if _container.cache == nil {
				_container.cache = map[string]*item{}
		}
		if _container.factories == nil {
				_container.factories = make([]func(args ...interface{}) interface{}, 2)
				_container.factories = _container.factories[:0]
		}
		return this
}

func (this *container) Get(id string, args ...interface{}) interface{} {
		this.Lock()
		defer this.UnLock()
		var v, ok = this.cache[id]
		if !ok {
				return ContainerNotFound
		}
		var (
				factory = this.factories[v.Index]
				ins     = factory(args...)
		)
		if ins != nil {
				v.Value = ins
		}
		return ins
}

func (this *container) GetInstance(name string, args ...interface{}) interface{} {
		this.Lock()
		defer this.UnLock()
		var v, ok = this.cache[name]
		if !ok {
				return ContainerNotFound
		}
		if v.Value != nil {
				return v.Value
		}
		var (
				factory = this.factories[v.Index]
				ins     = factory(args...)
		)
		if ins != nil {
				v.Value = ins
		}
		return ins
}

func (this *container) Register(name string, factory func(args ...interface{}) interface{}, force ...bool) ContainerInterface {
		if len(force) == 0 {
				force = append(force, false)
		}
		this.Lock()
		defer this.UnLock()
		if this.Exists(name) && !force[0] {
				return this
		}
		this.register(name, factory)
		return this
}

func (this *container) register(name string, factory func(args ...interface{}) interface{}) {
		this.factories = append(this.factories, factory)
		this.cache[name] = &item{
				Index: len(this.factories) - 1,
				Value: nil,
		}
}

func (this *container) Reset(name ...string) ContainerInterface {
		this.Lock()
		defer this.UnLock()
		if len(name) == 0 {
				this.cache = map[string]*item{}
				this.factories = this.factories[:0]
				return this
		}
		for _, id := range name {
				v, ok := this.cache[id]
				if !ok {
						continue
				}
				this.factories = append(this.factories[:v.Index], this.factories[v.Index:]...)
				delete(this.cache, id)
		}
		return this
}

func IsIocNotFound(v interface{}) bool {
		if e, ok := v.(error); ok {
				if e == ContainerNotFound {
						return true
				}
		}
		return false
}
