package services

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/libs"
		"reflect"
		"sync"
		"time"
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

func (this *BaseService) id(v interface{}) bson.ObjectId {
		if v == nil || v == "" {
				return ""
		}
		if str, ok := v.(string); ok {
				return bson.ObjectIdHex(str)
		}
		if id, ok := v.(bson.ObjectId); ok {
				return id
		}
		return ""
}

func (this *BaseService) isSet(m beego.M, key string) bool {
		if _, ok := m[key]; ok {
				return true
		}
		return false
}

func (this *BaseService) getObjectId(m beego.M, key string) bson.ObjectId {
		var v, ok = m[key]
		if !ok {
				return ""
		}
		if id, ok := v.(string); ok {
				return bson.ObjectIdHex(id)
		}
		if id, ok := v.(bson.ObjectId); ok {
				return id
		}
		return ""
}

func (this *BaseService) getStr(m beego.M, key string) string {
		var v, ok = m[key]
		if !ok {
				return ""
		}
		if value, ok := v.(string); ok {
				return value
		}
		if id, ok := v.(bson.ObjectId); ok {
				return id.Hex()
		}
		if str, ok := v.(fmt.Stringer); ok {
				return str.String()
		}
		return fmt.Sprintf("%v", v)
}

func (this *BaseService) getAny(m beego.M, key string) interface{} {
		var v, ok = m[key]
		if !ok {
				return nil
		}
		return v
}

func (this *BaseService) isEmpty(m beego.M, key string) bool {
		var v, ok = m[key]
		if !ok {
				return true
		}
		if v == nil || v == "" {
				return true
		}
		return false
}

func (this *BaseService) isZero(m beego.M, key string) bool {
		var v, ok = m[key]
		if !ok {
				return true
		}
		if v == nil || v == "" || v == 0 || v == false {
				return true
		}
		switch v.(type) {
		case beego.M:
				return len(v.(beego.M)) == 0
		case bson.M:
				return len(v.(bson.M)) == 0
		case map[string]interface{}:
				return len(v.(map[string]interface{})) == 0
		case map[interface{}]interface{}:
				return len(v.(map[interface{}]interface{})) == 0
		}
		return reflect.ValueOf(v).IsZero()
}

func (this *BaseService) getTime(query beego.M, key string) (time.Time, bool) {
		var v, ok = query[key]
		if !ok {
				return time.Now(), false
		}
		switch v.(type) {
		case string:
				str := v.(string)
				t, err := libs.GetTimeByNormalString(str)
				if err == nil {
						logs.Error(err)
						return time.Now(), false
				}
				return t, true
		case int64:
				t := v.(int64)
				return time.Unix(t, 0), true
		case time.Time:
				t := v.(time.Time)
				return t, true
		case *time.Time:
				t := v.(*time.Time)
				return *t, true
		}
		return time.Now(), false
}
