package libs

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config/env"
	"regexp"
	"strings"
)

// VariableParse 循环变量
func VariableParse(varStr string, i ...int) string {
	if len(i) <= 0 {
		i = append(i, 0)
	}
	if i[0] >= 3 {
		return varStr
	}
	var reg = regexp.MustCompile(`(\$\{[^\{\}]+\})?`)
	if !reg.Match([]byte(varStr)) {
		return varStr
	}
	var cache = make(map[string]struct {
		Value string
		Times int
	})
	vars := reg.FindAllString(varStr, -1)
	for _, it := range vars {
		def := ""
		key := strings.Replace(it, "${", "", 1)
		key = strings.Replace(key, "}", "", 1)
		if strings.Contains(key, "|") {
			arr := strings.SplitN(key, "|", 2)
			key = arr[0]
			def = arr[1]
		}
		// 是否有相同 key
		if v, ok := cache[key]; ok {
			if v.Times >= 1 && v.Value != "" {
				varStr = strings.ReplaceAll(varStr, it, v.Value)
				v.Times++
				cache[key] = v
				continue
			}
		}
		varN := env.Get(key, "<nil>")
		if varN == "<nil>" {
			varN = beego.AppConfig.String(key)
			if varN == "" && def == "" {
				continue
			}
			if varN == "" && def != "" {
				varN = def
			}
		}
		v, ok := cache[key]
		if ok && v.Times >= 1 {
			continue
		}
		if varN != "" {
			varN = VariableParse(varN, i[0]+1)
		}
		if varN != it {
			varStr = strings.ReplaceAll(varStr, it, varN)
		}
		if varN == "" && def != "" {
			cache[key] = struct {
				Value string
				Times int
			}{Value: "", Times: 1}
		} else {
			cache[key] = struct {
				Value string
				Times int
			}{Value: varN, Times: 1}
		}
	}
	return varStr
}
