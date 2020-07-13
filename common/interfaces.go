package common

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/astaxie/beego/session"
)

type RouterAction struct {
		Controller string
		Action     string
}

type BaseRequestContext interface {
		Method() string
		IsJsonStream() bool
		GetActionId() string
		GetJsonData() beego.M
		GetContentType() string
		GetJson() (beego.M, error)
		GetSession() session.Store
		GetParent() *beego.Controller
		GetInput() *context.BeegoInput
		JsonDecode(v interface{}) error
		GetInt(key string, def ...int) int
		GetControllerAction() *RouterAction
		GetBool(key string, def ...bool) bool
		GetString(key string, def ...string) string
		GetStrings(key string, def ...[]string) []string
		Session(key string, v ...interface{}) interface{}
		GetParam(key string, defaults ...interface{}) (interface{}, bool)
}
