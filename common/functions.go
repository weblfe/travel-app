package common

import (
		"fmt"
		"github.com/astaxie/beego"
		"os"
		"path"
)

func SmsDebugOn() {
		_ = os.Setenv("sms_debug_on", "1")
}

func SmsDebugOff() {
		_ = os.Setenv("sms_debug_on", "0")
}

func BasePath() string {
		return beego.AppPath
}

func StoragePath() string {
		if p := beego.AppConfig.String("storage_path"); p != "" {
				return p
		}
		if p := os.Getenv("STORAGE_PATH"); p != "" {
				return p
		}
		return path.Join(BasePath(), "/static/storage")
}

func Echo(v interface{}) {
		fmt.Printf("%v\n", v)
}

