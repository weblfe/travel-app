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

func Boolean(str string)bool  {
		if InArray(str,[]string{"true","True","Yes","yes","ok","Ok","On","on","1","TRUE"}) {
				return true
		}
		return false
}