package libs

import "regexp"

func IsCnMobile(mobile string) bool {
		var matcher = regexp.MustCompile(`^1[2-9][0-9]{9}$`)
		return matcher.Match([]byte(mobile))
}

func IsMobile(mobile string) bool {
		var matcher = regexp.MustCompile(`^\(\+[0-9]{2,4}\)[0-9]{11}$`)
		return matcher.Match([]byte(mobile))
}

func IsEmail(email string) bool {
		var matcher = regexp.MustCompile(`^\w{1,300}@\w{1,100}\.\w{2,10}$`)
		return matcher.Match([]byte(email))
}
