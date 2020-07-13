package transforms

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/services"
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
		return FilterUserAvatarUrl(m)
}

// 用户头像追加
func FilterUserAvatarUrl(m beego.M) beego.M {
		if url, ok := m["avatarUrl"]; ok && url != nil && url != "" {
				return m
		}
		if id, ok := m["avatarId"]; ok && id != nil && id != "" {
				m["avatarUrl"] = services.AvatarServerOf().GetAvatarUrlById(id.(string))
		} else {
				gender := m["gender"]
				if gender == nil {
						gender = 0
				}
				m["avatarUrl"] = services.AvatarServerOf().GetAvatarUrlDefault(gender.(int))
		}
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
