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
	var cache = make(map[string]int)
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
		if ok && v >= 1 {
			continue
		}
		if varN != "" {
			varN = VariableParse(varN, i[0]+1)
		}
		if varN != it {
			varStr = strings.ReplaceAll(varStr, it, varN)
		}
		cache[key] = 1
	}
	return varStr
}
