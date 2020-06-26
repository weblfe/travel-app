package libs

import (
		"crypto/md5"
		"encoding/hex"
		"github.com/astaxie/beego"
)

// md5
func Md5(str string) string {
		var ins = md5.New()
		ins.Write([]byte(str))
		return hex.EncodeToString(ins.Sum(nil))
}

// 加密密码
func PasswordHash(pass string, salt ...string) string {
		var (
				saltKey string
				ins     = md5.New()
		)
		if len(salt) == 0 {
				saltKey = beego.AppConfig.String("app_key")
				if saltKey == "" {
						saltKey = beego.AppConfig.String("appname")
				}
		} else {
				saltKey = salt[0]
		}
		ins.Write([]byte(pass + saltKey))
		return hex.EncodeToString(ins.Sum(nil))
}

// 密码验证
func PasswordVerify(encodePass string, pass string, salt ...string) bool {
		return encodePass == PasswordHash(pass, salt...)
}
