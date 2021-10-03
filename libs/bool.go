package libs

import "strings"

type EnumBool string

const (
	True  EnumBool = "true"
	Yes   EnumBool = "yes"
	Ok    EnumBool = "ok"
	On    EnumBool = "on"
	One   EnumBool = "1"
	No    EnumBool = "no"
	Not   EnumBool = "not"
	Off   EnumBool = "off"
	Zero  EnumBool = "0"
	False EnumBool = "false"
)

var (
	enumBoolMap = map[EnumBool]bool{
		True:  true,
		Yes:   true,
		Ok:    true,
		On:    true,
		One:   true,
		No:    false,
		Not:   false,
		Off:   false,
		Zero:  false,
		False: false,
	}
)

func (enum EnumBool) Bool() bool {
	var v, ok = enumBoolMap[enum.format()]
	if ok {
		return v
	}
	return false
}

func (enum EnumBool) Type() bool {
	var v = enum.format()
	if _, ok := enumBoolMap[v]; ok {
		return true
	}
	return false
}

func (enum EnumBool) format() EnumBool {
	var str = enum.String()
	return EnumBool(str)
}

func (enum EnumBool) String() string {
	return strings.ToLower(strings.TrimSpace(string(enum)))
}
