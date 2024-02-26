package controllers

import (
		"encoding/json"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/context"
		"github.com/astaxie/beego/session"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/middlewares"
		"net/http"
		"net/url"
		"reflect"
		"strconv"
		"strings"
)

type BaseController struct {
	beego.Controller
	_Request *beego.M
}

func (this *BaseController) Send(json common.ResponseJson) {
	this.Data["json"] = json
	this.ServeJSON()
}

func (this *BaseController) View(name string, data ...beego.M) {
	if len(data) > 0 {
		for _, m := range data {
			for k, v := range m {
				this.Data[k] = v
			}
		}
	}
	this.TplName = name
}

func (this *BaseController) GetParent() *beego.Controller {
	return &this.Controller
}

func (this *BaseController) GetInput() *context.BeegoInput {
	return this.Ctx.Input
}

func (this *BaseController) GetSession() session.Store {
	return this.Ctx.Input.CruSession
}

func (this *BaseController) GetContentType() string {
	return this.Ctx.Request.Header.Get("Content-Type")
}

func (this *BaseController) GetParam(key string, defaults ...interface{}) (interface{}, bool) {
	if len(defaults) == 0 {
		defaults = append(defaults, nil)
	}
	var v = this.Ctx.Input.Query(key)
	if v != "" {
		return v, true
	}
	if !this.IsJsonStream() {
		inputs := this.Input()
		if inputs == nil || len(inputs) == 0 {
			return defaults[0], false
		}
	}
	data, err := this.GetJson()
	if err != nil {
		return defaults[0], false
	}
	if v, ok := data[key]; ok {
		return v, ok
	}
	return defaults[0], false
}

func (this *BaseController) GetJson() (beego.M, error) {
	var data = make(beego.M)
	if this.IsJsonStream() {
		decoder := json.NewDecoder(this.Ctx.Request.Body)
		if err := decoder.Decode(&data); err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, fmt.Errorf("not json stream")
}

func (this *BaseController) IsJsonStream() bool {
	var typ = this.GetContentType()
	if strings.Contains(typ, "json") && !strings.Contains(typ, "jsonp") {
		return true
	}
	return false
}

func (this *BaseController) GetBody() []byte {
	return this.Ctx.Input.RequestBody
}

func (this *BaseController) JsonDecode(v interface{}) error {
	if !this.IsJsonStream() {
		return fmt.Errorf("is not json stream")
	}
	return json.Unmarshal(this.GetBody(), v)
}

func (this *BaseController) Session(key string, v ...interface{}) interface{} {
	if len(v) == 0 {
		return this.Ctx.Input.CruSession.Get(key)
	}
	return this.Ctx.Input.CruSession.Set(key, v)
}

func (this *BaseController) Method() string {
	return this.Ctx.Input.Method()
}

func (this *BaseController) GetControllerAction() *common.RouterAction {
	var data = new(common.RouterAction)
	data.Controller, data.Action = this.GetControllerAndAction()
	return data
}

func (this *BaseController) Cookie(key string, v ...interface{}) interface{} {
	var argc = len(v)
	if argc == 0 {
		return this.Ctx.GetCookie(key)
	}
	var value = v[0]
	// 删除
	if value == nil {
		this.Ctx.SetCookie(key, "")
		return true
	}
	if data, ok := value.(string); ok {
		if argc == 1 {
			this.Ctx.SetCookie(key, data)
			return true
		}
		if argc > 1 {
			this.Ctx.SetCookie(key, data, v[1:]...)
			return true
		}
	}

	return false
}

func (this *BaseController) GetHeader() http.Header {
	return this.Ctx.Request.Header
}

func (this *BaseController) SetHeader(key string, v string) {
	this.Ctx.Output.Header(key, v)
	return
}

func (this *BaseController) GetRequestContext() common.BaseRequestContext {
	return this
}

func (this *BaseController) GetString(key string, def ...string) string {
	if len(def) == 0 {
		def = append(def, "")
	}
	var v = this.Controller.GetString(key)
	if v != "" {
		return v
	}
	data := this.getJsonRequest()
	it, ok := data[key]
	if !ok {
		return def[0]
	}
	if str, ok := it.(string); ok {
		return str
	}
	return def[0]
}

func (this *BaseController) GetStrings(key string, def ...[]string) []string {
	if len(def) == 0 {
		def = append(def, []string{})
	}
	if this.Ctx.Request.Method == http.MethodGet {
		values, ok := this.query(key)
		if ok && values != nil {
			data := this.parseArray(values)
			if data != nil {
				return data
			}
		}
	}
	if this.IsJsonStream() {
		data := this.getJsonRequest()
		if len(data) == 0 {
			return def[0]
		}
		v, ok := data[key]
		if !ok {
			return def[0]
		}
		if arr, ok := v.([]string); ok {
			return arr
		}
		if arr, ok := v.([]interface{}); ok {
			var strArr []string
			for _, it := range arr {
				v, ok := it.(string)
				if !ok {
					continue
				}
				strArr = append(strArr, v)
			}
			return strArr
		}
		if str, ok := v.(string); ok {
			return strings.SplitN(str, ",", -1)
		}
	}
	return this.Controller.GetStrings(key, def...)
}

func (this *BaseController) GetInt(key string, def ...int) int {
	if len(def) == 0 {
		def = append(def, 0)
	}
	var v, err = this.Controller.GetInt(key, def...)
	if err == nil {
		return v
	}
	data := this.getJsonRequest()
	if v, ok := data[key]; ok {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Int:
			return v.(int)
		case reflect.Int8:
			return int(v.(int8))
		case reflect.Int64:
			return int(v.(int64))
		case reflect.Int16:
			return int(v.(int16))
		case reflect.Int32:
			return int(v.(int32))
		case reflect.String:
			n, err := strconv.Atoi(v.(string))
			if err == nil {
				return n
			}
		}
	}
	return def[0]
}

func (this *BaseController) GetBool(key string, def ...bool) bool {
	var v, err = this.Controller.GetBool(key, def...)
	if err == nil {
		return v
	}
	data := this.getJsonRequest()
	if v, ok := data[key]; ok {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			if v == 0 {
				return false
			}
			return true
		case reflect.String:
			if v == "1" || v == "True" || v == "true" {
				return true
			}
			return false
		case reflect.Bool:
			return v.(bool)
		}
	}
	return def[0]
}

func (this *BaseController) GetJsonData() beego.M {
	return this.getJsonRequest()
}

func (this *BaseController) GetActionId() string {
	return this.Ctx.Request.RequestURI
}

func (this *BaseController) getJsonRequest() beego.M {
	if this._Request == nil {
		this._Request = new(beego.M)
		if this.IsJsonStream() {
			_ = this.JsonDecode(this._Request)
		}
	}
	return *this._Request
}

func (this *BaseController) GetDriver() string {
	var driver = this.Ctx.Request.Header.Get("driver")
	if driver == "" {
		return "android"
	}
	return driver
}

func (this *BaseController) GetUserId() string {
	var sess session.Store
	if this.CruSession == nil {
		sess = this.Ctx.Input.CruSession
	}
	if sess == nil {
		return ""
	}
	var (
		userId = sess.Get(middlewares.AuthUserId)
	)
	if userId == nil || userId == "" {
		return ""
	}
	if id, ok := userId.(string); ok {
		return id
	}
	return ""
}

func (this *BaseController) parseArray(values interface{}) []string {
	switch values.(type) {
	case []string:
		return values.([]string)
	case string:
		v := values.(string)
		if strings.Contains(v, ",") {
			return strings.Split(v, ",")
		}
		var arr []string
		if err := json.Unmarshal([]byte(v), &arr); err == nil {
			return arr
		}
	case []byte:
		var arr []string
		if err := json.Unmarshal(values.([]byte), &arr); err == nil {
			return arr
		}
	}
	return nil
}

func (this *BaseController) query(key string, def ...interface{}) (interface{}, bool) {
	var query = this.Ctx.Request.RequestURI
	def = append(def, nil)
	if strings.Contains(query, "?") {
		values := strings.SplitN(query, "?", 2)
		if len(values) >= 2 {
			query = values[1]
		}
	}
	values, err := url.ParseQuery(query)
	if err != nil {
		return def[0], false
	}
	if v, ok := values[key]; ok {
		return v, true
	}
	if v, ok := values[key+"[]"]; ok {
		return v, true
	}
	return def[0], false
}
