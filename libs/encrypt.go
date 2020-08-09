package libs

import (
		"bytes"
		"crypto/des"
		"encoding/hex"
		"errors"
		"github.com/astaxie/beego/logs"
		"os"
		"strings"
)

const (
		_DefaultDesKey   = "travel@beego.2021.com"
		DesEnvKeyName    = "APP_ENCRYPT_DES_KEY"
		DesEncryptErrMsg = "need a multiple of the block size"
		DesDecryptErrMsg = "crypto/cipher: input not full blocks"
)

// 加密
func Encrypt(data string, key ...string) string {
		var (
				keyBit      = GetDesKey(key...)
				encode, err = DesEncrypt(data, keyBit)
		)
		if err == nil {
				return encode
		}
		logs.Error(err)
		return ""
}

// 解密
func Decrypt(data string, key ...string) string {
		var (
				keyBit      = GetDesKey(key...)
				decode, err = DesDecrypt(data, keyBit)
		)
		if err == nil {
				return decode
		}
		logs.Error(err)
		return ""
}

func GetDesKey(keys ...string) []byte {
		if len(keys) != 0 {
				return newDesKey(keys[0])
		}
		var key = os.Getenv(DesEnvKeyName)
		if key == "" {
				return newDesKey(_DefaultDesKey)
		}
		return newDesKey(key)
}

func newDesKey(key string) []byte {
		if key == "" {
				key = _DefaultDesKey
		}
		var (
				bits = []byte(key)
				size = len(bits)
		)
		if size < 8 {
				key = strings.Repeat(key, 8-size)
				bits = []byte(key)
				size = len(bits)
		}
		if size != 8 {
				bits = bits[:8]
		}
		return bits
}

func ZeroPadding(cipherText []byte, blockSize int) []byte {
		padding := blockSize - len(cipherText)%blockSize
		padText := bytes.Repeat([]byte{0}, padding)
		return append(cipherText, padText...)
}

func ZeroUnPadding(origData []byte) []byte {
		return bytes.TrimFunc(origData,
				func(r rune) bool {
						return r == rune(0)
				})
}

func DesEncrypt(text string, key []byte) (string, error) {
		src := []byte(text)
		block, err := des.NewCipher(key)
		if err != nil {
				return "", err
		}
		bs := block.BlockSize()
		src = ZeroPadding(src, bs)
		if len(src)%bs != 0 {
				return "", errors.New(DesEncryptErrMsg)
		}
		out := make([]byte, len(src))
		dst := out
		for len(src) > 0 {
				block.Encrypt(dst, src[:bs])
				src = src[bs:]
				dst = dst[bs:]
		}
		return hex.EncodeToString(out), nil
}

func DesDecrypt(decrypted string, key []byte) (string, error) {
		src, err := hex.DecodeString(decrypted)
		if err != nil {
				return "", err
		}
		block, err := des.NewCipher(key)
		if err != nil {
				return "", err
		}
		out := make([]byte, len(src))
		dst := out
		bs := block.BlockSize()
		if len(src)%bs != 0 {
				return "", errors.New(DesDecryptErrMsg)
		}
		for len(src) > 0 {
				block.Decrypt(dst, src[:bs])
				src = src[bs:]
				dst = dst[bs:]
		}
		out = ZeroUnPadding(out)
		return string(out), nil
}
