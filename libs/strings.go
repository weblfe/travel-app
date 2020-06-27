package libs

import "strings"

func MarkerMobile(mobile string) string {
		if !IsCnMobile(mobile) {
				return mobile
		}
		var arr = []rune(mobile)
		mobile = string(arr[0:3]) + strings.Repeat("*", 6) + string(arr[9:])
		return mobile
}
