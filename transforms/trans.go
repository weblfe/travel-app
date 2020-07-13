package transforms

import (
		"regexp"
		"strings"
)

func MarkerMobileTrans(mobile string) string {
		if !regexp.MustCompile(`^1[2-9][0-9]{9}$`).MatchString(mobile) {
				return mobile
		}
		var arr = []rune(mobile)
		mobile = string(arr[0:3]) + strings.Repeat("*", 6) + string(arr[9:])
		return mobile
}

