package libs

import "strings"

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
