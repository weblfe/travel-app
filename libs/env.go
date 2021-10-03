package libs

import (
	"os"
	"strconv"
)

func GetEnvBool(key string, def ...bool) bool {
	def = append(def, false)
	var value = os.Getenv(key)
	if value == "" {
		return def[0]
	}
	if v, err := strconv.ParseBool(value); err == nil {
		return v
	}
	var boolEnum = EnumBool(value)
	if boolEnum.Type() {
		return boolEnum.Bool()
	}
	return def[0]
}

func GetEnvOr(key string, def ...string) string {
	def = append(def, "")
	var value = os.Getenv(key)
	if value != "" {
		return value
	}
	return def[0]
}

func GetEnvInt(key string, def ...int) int {
	def = append(def, 0)
	var v = GetEnvOr(key)
	if v == "" {
		return def[0]
	}
	if n, err := strconv.Atoi(v); err == nil {
		return n
	}
	return def[0]
}

func GetEnvFloat(key string, def ...float64) float64 {
	def = append(def, 0)
	var v = GetEnvOr(key)
	if v == "" {
		return def[0]
	}
	if n, err := strconv.ParseFloat(v, 64); err == nil {
		return n
	}
	return def[0]
}
