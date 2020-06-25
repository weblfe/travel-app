package services

import (
		"sync"
)

type BaseService struct {
		container   *sync.Map
		ClassName   string
		Lock        sync.Mutex
		Attributes  map[string]interface{}
		Constructor func(args ...interface{}) interface{}
}

type Service interface {
		Class() string
		Hash() string
		Service() Service
		GetInstance(...interface{}) interface{}
		Invoker() func(args ...interface{}) interface{}
		GetAttribute(key string, defaults ...interface{}) interface{}
		SetAttribute(key string, value interface{}) Service
}

var (
		instanceLock sync.Once
		baseService  *BaseService
)

const (
		baseServiceClassName = "BaseService"
)

func newService() {
		baseService = new(BaseService)
		baseService.init()
}

func ServiceOf() *BaseService {
		if baseService == nil {
				instanceLock.Do(newService)
		}
		return baseService
}

func (this *BaseService) init() {
		if this.container == nil {
				this.container = new(sync.Map)
		}
		if this.Attributes == nil {
				this.Attributes = make(map[string]interface{})
		}
}

func (this *BaseService) Service() Service {
		return this
}

func (this *BaseService) GetInstance(args ...interface{}) interface{} {
		if len(args) == 0 {
				return this.get(this.Service().Class())
		}
		return this.resolve(args)
}

func (this *BaseService) get(name string) interface{} {
		if v, ok := this.container.Load(name); ok && v != nil {
				return v
		}
		return nil
}

func (this *BaseService) Class() string {
		return baseServiceClassName
}

func (this *BaseService) resolve(args []interface{}) Service {
		var invoker = this.Invoker()
		if invoker == nil {
				return nil
		}
		obj := invoker(args...)
		if service, ok := obj.(Service); ok {
				return service
		}
		return nil
}

func (this *BaseService) Invoker() func(args ...interface{}) interface{} {
		return this.Constructor
}

func (this *BaseService) Hash() string {
		return this.getHash()
}

func (this *BaseService) getHash() string {
		return ""
}

func (this *BaseService) GetAttribute(key string, defaults ...interface{}) interface{} {
		if v, ok := this.Attributes[key]; ok {
				return v
		}
		if len(defaults) == 0 {
				defaults = append(defaults, nil)
		}
		return defaults[0]
}

func (this *BaseService) SetAttribute(key string, value interface{}) Service {
		this.Attributes[key] = value
		return this
}
