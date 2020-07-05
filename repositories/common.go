package repositories

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"reflect"
		"time"
)

func isForbid(data *models.User) bool {
		return data.DeletedAt != 0 || data.Status != 1
}

// 获取多余字段
func filterUser(m beego.M) beego.M {
		delete(m, "deletedAt")
		delete(m, "passwordHash")
		if str, ok := m["mobile"]; ok && str != "" {
				m["mobile"] = libs.MarkerMobile(str.(string))
		}
		return m
}

// 过滤用户基础数据
func filterUserBase(m beego.M) beego.M {
		m = filterUser(m)
		delete(m, "accessTokens")
		delete(m, "registerWay")
		return filterUserAvatarUrl(m)
}

// 用户头像追加
func filterUserAvatarUrl(m beego.M) beego.M  {
		if id,ok:=m["avatarId"];ok &&id!=nil&&id!="" {
				m["avatarUrl"] = services.AvatarServerOf().GetAvatarUrlById(id.(string))
		}
		return m
}

// 过滤空数据字段
func filterEmpty(m beego.M) beego.M {
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
func filter(m beego.M, extras ...map[string]interface{}) beego.M {
		if len(extras) == 0 {
				extras = append(extras, map[string]interface{}{})
		}
		return filterEmptyMapper(filterEmpty(m))
}

// 过滤空mapper
func filterEmptyMapper(m beego.M) beego.M {
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
func filterAttachment(m beego.M) beego.M {
		if v, ok := m["id"]; ok {
				m["mediaId"] = v
		}
		return filterEmpty(m)
}

// 过滤器包装器
func FilterWrapper(filters ...func(m beego.M) beego.M) func(m beego.M) beego.M {
		if len(filters) <= 0 {
				return func(m beego.M) beego.M {
						return filter(m)
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
