package libs

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"log"
)

const (
	PasswordAPPKey = "71e920133ebb7d0a94b9daed8f6c2d9a"
)

// Md5 md5
func Md5(str string) string {
	var ins = md5.New()
	ins.Write([]byte(str))
	return hex.EncodeToString(ins.Sum(nil))
}

// PasswordHash 加密密码
func PasswordHash(pass string, salt ...string) string {
	var (
		saltKey string
		ins     = md5.New()
	)
	if len(salt) == 0 {
		salt = append(salt, beego.AppConfig.String("app_key"))
		if salt[0] == "" {
			salt = append(salt, PasswordAPPKey)
		}
	}
	saltKey = salt[0]
	ins.Write([]byte(pass + saltKey))
	log.Println(fmt.Sprintf("encode-password=%s", pass+saltKey))
	return hex.EncodeToString(ins.Sum(nil))
}

// PasswordVerify 密码验证
func PasswordVerify(encodePass string, pass string, salt ...string) bool {
	newPwd := PasswordHash(pass, salt...)
	log.Println(fmt.Sprintf("pass=%s,newPwd=%s,enCodePass=%s", pass, newPwd, encodePass))
	return encodePass == newPwd
}
