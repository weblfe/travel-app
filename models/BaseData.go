package models

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"reflect"
		"strings"
		"time"
)

type DataClass interface {
		Save() error
		InitDefault()
		Keys() []string
		Values() []interface{}
		Excludes(string) bool
		AddExcludeKeys(...string)
		RemoveExcludeKeys(...string)
		SetProvider(name string, v interface{})
		M(filter ...func(m beego.M) beego.M) beego.M
		AddFilters(filters ...func(m beego.M) beego.M)
		SetAttributes(m map[string]interface{}, safe ...bool)
}

type MapperAble interface {
		M(filters ...func(m beego.M) beego.M) beego.M
}

type dataClassImpl struct {
		data               *beego.M
		saver              func() error
		safeKeys           map[string]int
		filters            []func(m beego.M) beego.M
		attributesProvider func(map[string]interface{}, ...bool)
		dataProvider       func() beego.M
		defaultProvider    func()
}

const (
		SaverProvider      = "saverProvider"
		AttributesProvider = "attributesProvider"
		DataProvider       = "dataProvider"
		DefaultProvider    = "defaultProvider"
)

// 保存
func (this *dataClassImpl) Save() error {
		if this.saver != nil {
				return this.saver()
		}
		return common.NewErrors("miss saver", 2803)
}

// 是否排除键
func (this *dataClassImpl) Excludes(key string) bool {
		if this.safeKeys == nil || len(this.safeKeys) == 0 {
				return false
		}
		for k := range this.safeKeys {
				if k == key || strings.EqualFold(k, key) {
						return true
				}
		}
		return false
}

// 添加排除键
func (this *dataClassImpl) AddExcludeKeys(key ...string) {
		for _, k := range key {
				k = strings.ToLower(k)
				if _, ok := this.safeKeys[k]; ok {
						continue
				}
				this.safeKeys[k] = 1
		}
}

// 移除排除键
func (this *dataClassImpl) RemoveExcludeKeys(key ...string) {
		for _, k := range key {
				k = strings.ToLower(k)
				delete(this.safeKeys, k)
		}
}

// 获取数据对象
func (this *dataClassImpl) getData() beego.M {
		if this.data != nil {
				return *this.data
		}
		if this.dataProvider != nil {
				if d := this.dataProvider(); d != nil {
						this.data = &d
						return d
				}
		}
		return nil
}

// 过滤输出数据
func (this *dataClassImpl) M(filter ...func(m beego.M) beego.M) beego.M {
		var (
				data    = this.getData()
				filters = this.filters
		)
		filters = append(filters, filter...)
		for _, filter := range filters {
				if filter == nil {
						continue
				}
				data = filter(data)
		}
		return data
}

// 设置属性值
func (this *dataClassImpl) SetAttributes(m map[string]interface{}, safe ...bool) {
		if this.attributesProvider != nil {
				this.attributesProvider(m, safe...)
		}
}

// 设置服务提供函数
func (this *dataClassImpl) SetProvider(name string, v interface{}) {
		if reflect.TypeOf(v).Kind() != reflect.Func {
				return
		}
		switch name {
		case DataProvider:
				if fn, ok := v.(func() beego.M); ok && this.dataProvider == nil {
						this.dataProvider = fn
				}
		case SaverProvider:
				if fn, ok := v.(func() error); ok && this.saver == nil {
						this.saver = fn
				}
		case AttributesProvider:
				if fn, ok := v.(func(map[string]interface{}, ...bool)); ok && this.attributesProvider == nil {
						this.attributesProvider = fn
				}
		case DefaultProvider:
				if fn, ok := v.(func()); ok && this.defaultProvider == nil {
						this.defaultProvider = fn
				}
		}
}

// 设置objectId
func (this *dataClassImpl) SetObjectId(objId *bson.ObjectId, v interface{}) bool {
		if objId == nil || *objId != "" {
				return false
		}
		if id, ok := v.(bson.ObjectId); ok {
				*objId = id
				return true
		}
		if id, ok := v.(*bson.ObjectId); ok {
				*objId = *id
				return true
		}
		if id, ok := v.(string); ok {
				*objId = bson.ObjectIdHex(id)
				return true
		}
		return false
}

// 设置时间
func (this *dataClassImpl) SetTime(tObj *time.Time, v interface{}, force ...bool) bool {
		if tObj == nil {
				return false
		}
		if force[0] || tObj.IsZero() {
				if t, ok := v.(time.Time); ok && !t.IsZero() {
						*tObj = t
						return true
				}
				if t, ok := v.(*time.Time); ok && !t.IsZero() {
						*tObj = *t
						return true
				}
				if t, ok := v.(int64); ok && t > 0 {
						tt := time.Unix(t, 0)
						if tt.IsZero() {
								return false
						}
						*tObj = tt
						return true
				}
		}
		return false
}

// 设置mapper
func (this *dataClassImpl) SetMapper(tObj *beego.M, v interface{}, force ...bool) bool {
		if tObj == nil {
				return false
		}
		if force[0] || len(*tObj) == 0 {
				switch v.(type) {
				case beego.M:
						*tObj = v.(beego.M)
						return true
				case *beego.M:
						*tObj = *v.(*beego.M)
						return true
				case bson.M:
						*tObj = beego.M(v.(bson.M))
						return true
				case *bson.M:
						*tObj = beego.M(*v.(*bson.M))
						return true
				case map[string]interface{}:
						*tObj = v.(map[string]interface{})
						return true
				case *map[string]interface{}:
						*tObj = *v.(*map[string]interface{})
						return true
				}
				if fn, ok := v.(func() beego.M); ok {
						*tObj = fn()
						return true
				}
				if m, ok := v.(MapperAble); ok {
						*tObj = m.M()
						return true
				}
		}
		return false
}

// 数据键集合
func (this *dataClassImpl) Keys() []string {
		var keys []string
		for k, _ := range this.getData() {
				keys = append(keys, k)
		}
		return keys
}

// 数据值集合
func (this *dataClassImpl) Values() []interface{} {
		var values []interface{}
		for _, v := range this.getData() {
				values = append(values, v)
		}
		return values
}

// 遍历
func (this *dataClassImpl) Foreach(each func(k string, v interface{}) bool) {
		for k, v := range this.getData() {
				if !each(k, v) {
						break
				}
		}
}

func (this *dataClassImpl) IsEmpty(v interface{}) bool {
		// 空
		if v == nil || v == "" || v == 0 {
				return true
		}
		// 空时间
		if t, ok := v.(time.Time); ok && t.IsZero() {
				return true
		}
		return false
}

// 遍历
func (this *dataClassImpl) Map(each func(k string, v interface{}, result interface{}) interface{}, result ...interface{}) interface{} {
		var ret = result[0]
		for k, v := range this.getData() {
				ret = each(k, v, ret)
		}
		return ret
}

// 合并
func (this *dataClassImpl) Merger(m beego.M, m2 beego.M) beego.M {
		return Merger(m, m2)
}

// 设置字符串
func (this *dataClassImpl) SetString(str *string, v interface{}) bool {
		if v == nil || str == nil {
				return false
		}
		switch v.(type) {
		case string:
				*str = v.(string)
				return true
		case *string:
				*str = *v.(*string)
				return true
		case []byte:
				*str = string(v.([]byte))
				return true
		case []string:
				*str = strings.Join(v.([]string), ",")
				return true
		}
		if s, ok := v.(bson.ObjectId); ok {
				*str = s.Hex()
				return true
		}
		if s, ok := v.(fmt.Stringer); ok {
				*str = s.String()
				return true
		}
		if s, ok := v.(func() string); ok {
				*str = s()
				return true
		}
		*str = fmt.Sprintf("%v", v)
		return true
}

// 相同类型赋值
func (this *dataClassImpl) SetSameTypeValue(obj interface{}, v interface{}) bool {
		if obj == nil || v == nil {
				return false
		}
		var (
				tObj, tV = reflect.TypeOf(obj), reflect.TypeOf(v)
		)

		if tObj.Elem().Kind() == tV.Kind() && tObj.Elem().Name() == tV.Name() {
				reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(v))
				return true
		}
		return false
}

// 添加过滤器
func (this *dataClassImpl) AddFilters(filters ...func(m beego.M) beego.M) {
		if this.filters == nil {
				this.filters = make([]func(m beego.M) beego.M, 2)
		}
		this.filters = append(this.filters, filters...)
}

// 初始化默认值
func (this *dataClassImpl) InitDefault() {
		if this.defaultProvider == nil {
				return
		}
		this.defaultProvider()
		return
}

func Merger(m beego.M, m2 beego.M) beego.M {
		for k, v := range m2 {
				m[k] = v
		}
		return m
}