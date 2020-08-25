package models

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"reflect"
		"strconv"
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

func (this *dataClassImpl) GetNow() time.Time {
		return time.Now().Local()
}

func (this *dataClassImpl) GetId() bson.ObjectId {
		return bson.NewObjectId()
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

func (this *dataClassImpl) reset() {
		this.data = nil
		return
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

// 时间戳格式
func (this *dataClassImpl) GetFormatterTime(key string) func(data beego.M) beego.M {
		return func(data beego.M) beego.M {
				var value, ok = data[key]
				if !ok {
						return data
				}
				if v, ok := value.(int64); ok {
						data[key] = time.Unix(v, 0)
						return data
				}
				if v, ok := value.(int); ok {
						data[key] = time.Unix(int64(v), 0)
						return data
				}
				return data
		}
}

// 字段过滤器
func (this *dataClassImpl) GetKeysFilter(keys []string, excludes ...bool) func(data beego.M) beego.M {
		if len(excludes) == 0 {
				excludes = append(excludes, true)
		}
		return func(data beego.M) beego.M {
				var exclude = excludes[0]
				if len(keys) == 0 {
						return data
				}
				if !exclude {
						var results = beego.M{}
						for _, key := range keys {
								if v, ok := data[key]; ok {
										results[key] = v
								}
						}
						return results
				}
				for _, key := range keys {
						delete(data, key)
				}
				return data
		}
}

// 字段转换器
func (this *dataClassImpl) GetTransformFilterByKey(key string, trans func(v interface{}) interface{}) func(data beego.M) beego.M {
		return func(data beego.M) beego.M {
				if v, ok := data[key]; ok {
						data[key] = trans(v)
				}
				return data
		}
}

// 字段转换器
func (this *dataClassImpl) GetTransformFilter(transform func(key string, v interface{}, data *beego.M)) func(data beego.M) beego.M {
		return func(data beego.M) beego.M {
				for key, v := range data {
						transform(key, v, &data)
				}
				return data
		}
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

// bson Map 赋值
func (this *dataClassImpl) SetBsonMapper(tObj *bson.M, v interface{}, force ...bool) bool {
		return this.SetMapper((*beego.M)(tObj), v, force...)
}

//  Map 赋值
func (this *dataClassImpl) SetMap(tObj *map[string]interface{}, v interface{}, force ...bool) bool {
		return this.SetMapper((*beego.M)(tObj), v, force...)
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

// 设置数字
func (this *dataClassImpl) SetNumInt(num *int, v interface{}) bool {
		switch v.(type) {
		case int:
				*num = v.(int)
		case int8:
				*num = int(v.(int8))
		case int16:
				*num = int(v.(int16))
		case int32:
				*num = int(v.(int32))
		case int64:
				*num = int(v.(int64))
		case float64:
				*num = int(v.(float64))
		case float32:
				*num = int(v.(float32))
		case string:
				if n, err := strconv.Atoi(v.(string)); err == nil {
						*num = n
						return true
				}
				return false
		default:
				return false
		}
		return false
}

// 设置数字
func (this *dataClassImpl) SetNumIntN(num *int64, v interface{}) bool {
		switch v.(type) {
		case int:
				*num = int64(v.(int))
		case int8:
				*num = int64(v.(int8))
		case int16:
				*num = int64(v.(int16))
		case int32:
				*num = int64(v.(int32))
		case int64:
				*num = v.(int64)
		case float64:
				*num = int64(v.(float64))
		case float32:
				*num = int64(v.(float32))
		case string:
				if n, err := strconv.Atoi(v.(string)); err == nil {
						*num = int64(n)
						return true
				}
				return false
		case time.Time:
				var t = v.(time.Time)
				*num = t.Unix()
				return true
		default:
				return false
		}
		return false
}

// 设置数字
func (this *dataClassImpl) SetBool(value *bool, v interface{}) bool {
		switch v.(type) {
		case int:
				n := v.(int)
				if n > 0 {
						*value = true
				} else {
						*value = false
				}
		case int8:
				n := v.(int8)
				if n > 0 {
						*value = true
				} else {
						*value = false
				}
		case int16:
				n := v.(int16)
				if n > 0 {
						*value = true
				} else {
						*value = false
				}
		case int32:
				n := v.(int32)
				if n > 0 {
						*value = true
				} else {
						*value = false
				}
		case int64:
				n := v.(int64)
				if n > 0 {
						*value = true
				} else {
						*value = false
				}
		case float64:
				n := v.(float64)
				if n > 0 {
						*value = true
				} else {
						*value = false
				}
		case float32:
				n := v.(float32)
				if n > 0 {
						*value = true
				} else {
						*value = false
				}
		case string:
				n := v.(string)
				if n == "" {
						*value = false
						return true
				}
				for _, it := range []string{"yes", "ok", "on", "1", "true", "open"} {
						if it == n || strings.EqualFold(it, n) {
								*value = true
								return true
						}
				}
		case bool:
				*value = v.(bool)
		default:
				return false
		}
		return false
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
				this.filters = this.filters[:0]
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

func (this *dataClassImpl) SetStringArr(arr *[]string, v interface{}) bool {
		if str, ok := v.(string); ok && str != "" {
				*arr = strings.SplitN(str, ",", -1)
				return true
		}
		if it, ok := v.([]string); ok && len(it) > 0 {
				*arr = it
				return true
		}
		return false
}

func Merger(m beego.M, m2 beego.M) beego.M {
		for k, v := range m2 {
				m[k] = v
		}
		return m
}

// 是否数字类型
func IsNumber(v interface{}) bool {
		switch v.(type) {
		case int:
				return true
		case int8:
				return true
		case int16:
				return true
		case int32:
				return true
		case int64:
				return true
		case float64:
				return true
		case float32:
				return true
		}
		return false
}

// 是否为空
func IsEmpty(v interface{}) bool {
		var getValue = reflect.ValueOf(v)
		return getValue.IsZero() || getValue.IsNil()
}
