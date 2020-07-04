package libs

import "github.com/astaxie/beego"

func MapMerge(m ...beego.M) beego.M {
		if len(m) == 0 {
				return beego.M{}
		}
		var result = m[0]
		for i := 1; i < len(m); i++ {
				for k, v := range m[i] {
						result[k] = v
				}
		}
		return result
}

func Merge(m...map[string]interface{}) map[string]interface{}  {
		if len(m) == 0 {
				return map[string]interface{}{}
		}
		var result = m[0]
		for i := 1; i < len(m); i++ {
				for k, v := range m[i] {
						result[k] = v
				}
		}
		return result
}
