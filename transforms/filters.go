package transforms

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"reflect"
		"time"
)

// 获取多余字段
func FilterUser(m beego.M) beego.M {
		delete(m, "deletedAt")
		delete(m, "passwordHash")
		if str, ok := m["mobile"]; ok && str != "" {
				m["mobile"] = MarkerMobileTrans(str.(string))
		}
		return m
}

// 过滤用户基础数据
func FilterUserBase(m beego.M) beego.M {
		m = FilterUser(m)
		delete(m, "accessTokens")
		delete(m, "registerWay")
		return m
}

// 过滤空数据字段
func FilterEmpty(m beego.M) beego.M {
		for k, v := range m {
				if v == "" || v == nil {
						delete(m, k)
						continue
				}
				getValue := reflect.ValueOf(v)
				if getValue.IsZero() {
						delete(m, k)
						continue
				}
				if t, ok := v.(time.Time); ok {
						if t.IsZero() {
								delete(m, k)
						}
				}

		}
		return m
}

// 过滤空数据
func FilterEmptyWithOutNumber(m beego.M) beego.M {
		for k, v := range m {
				if v == "" || v == nil {
						delete(m, k)
						continue
				}
				getValue := reflect.ValueOf(v)
				if getValue.IsZero() && !IsNumber(v) {
						delete(m, k)
						continue
				}
				if t, ok := v.(time.Time); ok {
						if t.IsZero() {
								delete(m, k)
						}
				}

		}
		return m
}

func IsNumber(value interface{}) bool {
		var v = reflect.TypeOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				return true
		case reflect.Float32, reflect.Float64:
				return true
		case reflect.Complex64, reflect.Complex128:
				return true
		}
		return false
}

// 过滤 任意空值
func Filter(m beego.M, extras ...map[string]interface{}) beego.M {
		if len(extras) == 0 {
				extras = append(extras, map[string]interface{}{})
		}
		return FilterEmptyMapper(FilterEmpty(m))
}

// 过滤空mapper
func FilterEmptyMapper(m beego.M) beego.M {
		for key, v := range m {
				obj, ok := v.(map[string]interface{})
				if ok && len(obj) == 0 {
						delete(m, key)
						continue
				}
				getValue := reflect.ValueOf(v)
				if getValue.Kind() == reflect.Map && getValue.Len() == 0 {
						delete(m, key)
						continue
				}
				obj, ok = v.(beego.M)
				if ok && len(obj) == 0 {
						delete(m, key)
						continue
				}
				obj, ok = v.(bson.M)
				if ok && len(obj) == 0 {
						delete(m, key)
				}
		}
		return m
}

// attachment 过滤
func FilterAttachment(m beego.M) beego.M {
		if v, ok := m["id"]; ok {
				m["mediaId"] = v
		}
		return FilterEmpty(m)
}

// 过滤器包装器
func FilterWrapper(filters ...func(m beego.M) beego.M) func(m beego.M) beego.M {
		if len(filters) <= 0 {
				return func(m beego.M) beego.M {
						return Filter(m)
				}
		}
		return func(m beego.M) beego.M {
				for _, filter := range filters {
						m = filter(m)
				}
				return m
		}
}

// 字段过滤器
// files []string 字段
// exclude bool # true : 过滤对应字段, false: 保留对应字段
func FieldsFilter(fields []string, exclude ...bool) func(m beego.M) beego.M {
		if len(exclude) == 0 {
				exclude = append(exclude, true)
		}
		return func(m beego.M) beego.M {
				for key, _ := range m {
						for _, k := range fields {
								if k != key {
										// 保存
										if !exclude[0] {
												continue
										}
										continue
								}
								// 排除
								delete(m, key)
						}
				}
				return m
		}
}

// 时间转时间戳
func FilterTimeToInt64(m beego.M, keys ...string) beego.M {
		if len(keys) == 0 {
				keys = append(keys, "createdAt", "updatedAt", "deletedAt")
		}
		for _, key := range keys {
				v, ok := m[key]
				if !ok {
						continue
				}
				switch v.(type) {
				case time.Time:
						t := v.(time.Time)
						if t.IsZero() {
								m[key] = 0
								break
						}
						m[key] = t.Unix()
				case *time.Time:
						t := v.(time.Time)
						if t.IsZero() {
								m[key] = 0
								break
						}
						m[key] = t.Unix()
				case int64:
						break
				case string:
						if v == "" {
								break
						}
						t, err := time.Parse(time.RFC3339, v.(string))
						if err == nil {
								m[key] = t.Unix()
						}
				}
		}
		return m
}
