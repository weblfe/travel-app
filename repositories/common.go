package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"reflect"
		"time"
)

func isForbid(data *models.User) bool {
		return data.DeletedAt != 0 || data.Status != 1
}

// 获取多余字段
func filterUser(m beego.M) beego.M {
		delete(m, "deleted_at")
		delete(m, "password")
		if str, ok := m["mobile"]; ok && str != "" {
				m["mobile"] = libs.MarkerMobile(str.(string))
		}
		return m
}

// 过滤用户基础数据
func filterUserBase(m beego.M) beego.M {
		m = filterUser(m)
		delete(m, "access_tokens")
		return m
}

// 过滤空数据字段
func filterEmpty(m beego.M, number ...bool) beego.M {
		if len(number) == 0 {
				number = append(number, true)
		}
		for k, v := range m {
				if v == "" || v == nil {
						delete(m, k)
						continue
				}
				if v == 0 && number[0] {
						delete(m, k)
						continue
				}
				ty := reflect.TypeOf(v)
				if ty.Kind() == reflect.Slice || ty.Kind() == reflect.Array {
						if ty.Len() <= 0 {
								delete(m, k)
						}
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

func filterAttachment(m beego.M) beego.M {
		if v, ok := m["id"]; ok {
				m["mediaId"] = v
		}
		return filterEmpty(m)
}