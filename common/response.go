package common

import (
		"encoding/json"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/config/env"
		"github.com/globalsign/mgo/bson"
		"strings"
		"sync"
)

type ResponseJson interface {
		fmt.Stringer
		GetCode() int
		GetMsg() string
		GetError() Errors
		GetData() interface{}
		UnJson(bytes []byte) error
		Json() ([]byte, error)
		Set(key string, v interface{}) ResponseJson
		Has(key string) bool
		Empty() bool
		IsSuccess() bool
		Get(key string, defaults ...interface{}) interface{}
		GetDataByKey(key string, defaults ...interface{}) interface{}
}

type ResponseImpl struct {
		Code    int         `json:"code"`
		Message string      `json:"msg"`
		Error   Errors      `json:"error"`
		Data    interface{} `json:"data"`
}

type GetterEntry interface {
		Get(key string, defaults ...interface{}) interface{}
}

type GetEntry interface {
		Get(key string) (interface{}, bool)
}

type SetEntry interface {
		Set(key string, v interface{}) bool
}

type SetterEntry interface {
		Set(key string, v interface{})
}

func (this *ResponseImpl) UnJson(bytes []byte) error {
		return json.Unmarshal(bytes, this)
}

func (this *ResponseImpl) Json() ([]byte, error) {
		return json.Marshal(this)
}

func (this *ResponseImpl) GetCode() int {
		return this.Code
}

func (this *ResponseImpl) GetMsg() string {
		return this.Message
}

func (this *ResponseImpl) GetError() Errors {
		return this.Error
}

func (this *ResponseImpl) GetData() interface{} {
		return this.Data
}

func (this *ResponseImpl) Get(key string, defaults ...interface{}) interface{} {
		switch key {
		case "data":
				fallthrough
		case "Data":
				return this.GetData()
		case "Msg":
				fallthrough
		case "msg":
				fallthrough
		case "Message":
				fallthrough
		case "message":
				return this.GetMsg()
		case "code":
				fallthrough
		case "Code":
				return this.GetCode()
		case "err":
				fallthrough
		case "error":
				return this.GetError()
		}
		if strings.Index(key, "data.") == 0 {
				key = strings.Replace(key, "data.", "", 1)
				return this.GetDataByKey(key, defaults)
		}
		return nil
}

func (this *ResponseImpl) GetDataByKey(key string, defaults ...interface{}) interface{} {
		if len(defaults) == 0 {
				defaults = append(defaults, nil)
		}
		if get, ok := this.Data.(GetEntry); ok {
				if v, b := get.Get(key); b {
						return v
				}
				return defaults[0]
		}
		if get, ok := this.Data.(GetterEntry); ok {
				return get.Get(key, defaults...)
		}
		if mp, ok := this.Data.(map[string]interface{}); ok {
				if v, ok := mp[key]; ok {
						return v
				}
		}
		if mp, ok := this.Data.(bson.M); ok {
				if v, ok := mp[key]; ok {
						return v
				}
		}
		if mp, ok := this.Data.(beego.M); ok {
				if v, ok := mp[key]; ok {
						return v
				}
		}
		if mp, ok := this.Data.(*map[string]interface{}); ok {
				if v, ok := (*mp)[key]; ok {
						return v
				}
		}
		if mp, ok := this.Data.(sync.Map); ok {
				if v, ok := mp.Load(key); ok {
						return v
				}
		}
		if mp, ok := this.Data.(*sync.Map); ok {
				if v, ok := mp.Load(key); ok {
						return v
				}
		}
		return defaults[0]
}

func (this *ResponseImpl) Set(key string, v interface{}) ResponseJson {
		switch key {
		case "data":
				fallthrough
		case "Data":
				this.Data = v
				return this
		case "err":
				fallthrough
		case "Error":
				fallthrough
		case "error":
				if err, ok := v.(Errors); ok {
						this.Error = err
				}
				return this
		case "msg":
				fallthrough
		case "Msg":
				fallthrough
		case "message":
				fallthrough
		case "Message":
				if msg, ok := v.(string); ok {
						this.Message = msg
				}
				return this
		case "code":
				fallthrough
		case "Code":
				if code, ok := v.(int); ok {
						this.Code = code
				}
				return this
		}
		if strings.Index(key, "data.") == 0 {
				key = strings.Replace(key, "data.", "", 1)
				return this.SetDataByKey(key, v)
		}
		return this
}

func (this *ResponseImpl) SetDataByKey(key string, v interface{}) ResponseJson {
		if set, ok := this.Data.(SetEntry); ok {
				set.Set(key, v)
				return this
		}
		if set, ok := this.Data.(SetterEntry); ok {
				set.Set(key, v)
				return this
		}
		if mp, ok := this.Data.(map[string]interface{}); ok {
				mp[key] = v
				return this
		}
		if mp, ok := this.Data.(*map[string]interface{}); ok {
				(*mp)[key] = v
				return this
		}
		if mp, ok := this.Data.(sync.Map); ok {
				mp.Store(key, v)
				return this
		}
		if mp, ok := this.Data.(*sync.Map); ok {
				mp.Store(key, v)
				return this
		}
		return this
}

func (this *ResponseImpl) Has(key string) bool {
		switch key {
		case "data":
				fallthrough
		case "Data":
				return this.Data != nil
		case "err":
				fallthrough
		case "Error":
				fallthrough
		case "error":
				return this.Error != nil
		case "msg":
				fallthrough
		case "Msg":
				fallthrough
		case "message":
				fallthrough
		case "Message":
				return this.Message != ""
		case "code":
				fallthrough
		case "Code":
				return this.Code != -1
		}
		if strings.Index(key, "data.") == 0 {
				key = strings.Replace(key, "data.", "", 1)
				if v := this.GetDataByKey(key, nil); v == nil {
						return false
				}
				return true
		}
		return false
}

func (this *ResponseImpl) Empty() bool {
		if this.Data == nil && this.Error == nil && this.Message == "" && this.Code == -1 {
				return true
		}
		return false
}

func (this *ResponseImpl) IsSuccess() bool {
		if this.Empty() {
				return false
		}
		return this.Code == SuccessCode
}

func (this *ResponseImpl) init(args ...interface{}) {
		for _, arg := range args {
				if v, ok := arg.(struct{}); ok && this.Data == nil {
						this.Data = v
						continue
				}
				if v, ok := arg.(map[string]interface{}); ok && this.Data == nil {
						this.Data = v
						continue
				}
				if v, ok := arg.(beego.M); ok && this.Data == nil {
						this.Data = v
						continue
				}
				if v, ok := arg.(GetterEntry); ok && this.Data == nil {
						this.Data = v
						continue
				}
				if v, ok := arg.(GetEntry); ok && this.Data == nil {
						this.Data = v
						continue
				}
				if v, ok := arg.(SetterEntry); ok && this.Data == nil {
						this.Data = v
						continue
				}
				if v, ok := arg.(SetEntry); ok && this.Data == nil {
						this.Data = v
						continue
				}
				if v, ok := arg.(int); ok && this.Code == -1 {
						this.Code = v
						continue
				}
				if v, ok := arg.(string); ok && this.Message == "" {
						this.Message = v
						continue
				}
				if v, ok := arg.(Errors); ok && this.Error == nil {
						this.Error = v
				}
		}
}

func (this *ResponseImpl) String() string {
		if data, err := json.Marshal(this); err == nil {
				return string(data)
		}
		return fmt.Sprintf(
				`{"code":%d,"msg":"%s","data":"%s","error":"%s"}`,
				this.GetCode(), this.GetMsg(), this.GetData(), this.GetError(),
		)
}

func NewResponse(args ...interface{}) ResponseJson {
		var resp = new(ResponseImpl)
		resp.Code = -1
		resp.init(args...)
		return resp
}

// 构造成功响应体
// data interface
// msg  string
// code int
func NewSuccessResp(data interface{}, args ...interface{}) ResponseJson {
		if len(args) == 0 {
				args = append(args, Success)
		}
		if len(args) < 2 {
				args = append(args, SuccessCode)
		}
		var res = NewResponse(args...)
		res.Set("data", data)
		if !res.Has("code") {
				res.Set("code", args[0])
		}
		if !res.Has("msg") {
				res.Set("msg", args[1])
		}
		return res
}

// 未登陆响应
// data interface
// msg  string
// code int
func NewUnLoginResp(data ...interface{}) ResponseJson {
		var resp = NewResponse(data...)
		resp.Set("code", UnLoginCode)
		if !resp.Has("msg") {
				resp.Set("msg", UnLoginError)
		}
		return resp
}

// 访问控制
// data interface
// msg  string
// code int
func NewAccessLimitResp(data ...interface{}) ResponseJson {
		var resp = NewResponse(data...)
		resp.Set("code", LimitCode)
		if !resp.Has("msg") {
				resp.Set("msg", LimitError)
		}
		return resp
}

// 请求参数异常
// data interface
// msg  string
// code int
func NewInvalidParametersResp(data ...interface{}) ResponseJson {
		var resp = NewResponse(data...)
		if !resp.Has("code") {
				resp.Set("code", InvalidParametersCode)
		}
		if !resp.Has("msg") {
				resp.Set("msg", InvalidParametersError)
		}
		return resp
}

// 请求参数异常
// err Errors
// msg  string
// code int
// data interface{}
func NewErrorResp(err Errors, args ...interface{}) ResponseJson {
		var resp = NewResponse(args...)
		resp.Set("err", err)
		if !resp.Has("code") {
				resp.Set("code", ErrorCode)
		}
		if !resp.Has("msg") {
				resp.Set("msg", Error)
		}
		return resp
}

// 构造成功响应体
// data interface
// msg  string
// code int
func NewFailedResp(code int, args ...interface{}) ResponseJson {
		args = append(args, code)
		if len(args) < 2 {
				args = append(args, ServiceFailedError)
		}
		var res = NewResponse(args...)
		if !res.Has("code") {
				res.Set("code", args[0])
		}
		if !res.Has("msg") {
				res.Set("msg", args[1])
		}
		if !res.Has("error") {
				res.Set("error",NewErrors(res.Get("code"),res.Get("msg")))
		}
		return res
}

// 请求参数异常
// err Errors
// msg  string
// code int
// data interface{}
func NewInDevResp(api string, args ...interface{}) ResponseJson {
		var resp = NewResponse(args...)
		server := env.Get("SERVER_DOMAIN", "")
		errMsg := server + " api: " + api
		resp.Set("err", NewErrors(errMsg, DevelopCode))
		if !resp.Has("code") {
				resp.Set("code", DevelopCode)
		}
		if !resp.Has("msg") {
				resp.Set("msg", DevelopCodeError)
		}
		return resp
}

// 权限异常
// err Errors
// msg  string
// code int
// data interface{}
func NewPermissionResp(args ...interface{}) ResponseJson {
		var resp = NewResponse(args...)
		if !resp.Has("code") {
				resp.Set("code", PermissionCode)
		}
		if !resp.Has("msg") {
				resp.Set("msg", PermissionError)
		}
		return resp
}
