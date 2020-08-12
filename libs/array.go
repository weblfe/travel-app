package libs

import (
		"strings"
		"time"
)

func InArray(value string, array []string, fold ...bool) bool {
		if len(fold) <= 0 {
				fold = append(fold, false)
		}
		for _, it := range array {
				if fold[0] {
						if strings.EqualFold(value, it) {
								return true
						}
						continue
				}
				if value == it {
						return true
				}
		}
		return false
}

// 创建hash
func NewHashMapper(v ...interface{}) []interface{} {
		if len(v) == 0 {
				return nil
		}
		var arr []interface{}
		for _, it := range v {
				arr = append(arr, it)
		}
		return arr
}

// 创建hash
func NewIntegerArray(v ...interface{}) []int {
		var arr []int
		for _, n := range v {
				num, ok := Integer(n)
				if !ok {
						continue
				}
				arr = append(arr, num)
		}
		return arr
}

// 整型
func Integer(v interface{}) (int, bool) {
		switch v.(type) {
		case int:
				return v.(int), true
		case int64:
				return int(v.(int64)), true
		case int32:
				return int(v.(int32)), true
		case int16:
				return int(v.(int16)), true
		case int8:
				return int(v.(int8)), true
		case time.Time:
				t := v.(time.Time)
				return int(t.Unix()), true
		}
		return 0, false
}
