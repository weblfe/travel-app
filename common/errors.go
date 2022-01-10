package common

import "fmt"

type Errors interface {
		Code() int
		Msg() string
		Parent() Errors
		Set(key string, v interface{}) Errors
		error
}

type ErrorImpl struct {
		ErrCode int    `json:"errno"`
		ErrMsg  string `json:"errmsg"`
		parent  Errors
}

// NewErrors ErrCode string
// ErrMsg  int
func NewErrors(args ...interface{}) Errors {
		var err = new(ErrorImpl)
		err.ErrCode = -1
		err.init(args...)
		return err
}

// ErrorWrap err Errors
// ErrCode string
// ErrMsg  int
func ErrorWrap(err Errors, args ...interface{}) Errors {
		if len(args) == 0 {
				args = []interface{}{}
		}
		args = append([]interface{}{err}, args...)
		return NewErrors(args...)
}

func (this *ErrorImpl) init(args ...interface{}) {
		for _, arg := range args {
				if err, ok := arg.(Errors); ok && this.parent == nil {
						this.parent = err
						this.ErrCode = err.Code()
						this.ErrMsg = err.Msg()
				}
				if msg, ok := arg.(string); ok && this.ErrMsg == "" {
						this.ErrMsg = msg
				}
				if msg, ok := arg.(error); ok && this.ErrMsg == "" {
						this.ErrMsg = msg.Error()
				}
				if code, ok := arg.(int); ok && this.ErrCode == -1 {
						this.ErrCode = code
				}
		}
}

func (this *ErrorImpl) Error() string {
		return fmt.Sprintf(`{"errno":%d,"errmsg":"%s"}`, this.Code(), this.Msg())
}

func (this *ErrorImpl) String() string {
		return this.Error()
}

func (this *ErrorImpl) Parent() Errors {
		return this.parent
}

func (this *ErrorImpl) Code() int {
		if this.ErrCode == 0 && this.parent != nil {
				return this.parent.Code()
		}
		return this.ErrCode
}

func (this *ErrorImpl) Msg() string {
		if this.ErrMsg == "" && this.parent != nil {
				return this.parent.Msg()
		}
		return this.ErrMsg
}

func (this *ErrorImpl) Set(key string, v interface{}) Errors {
		switch key {
		case "ErrCode":
				fallthrough
		case "errno":
				this.ErrCode = v.(int)
		case "ErrMsg":
		case "errmsg":
				this.ErrMsg = v.(string)
		case "parent":
				this.parent = v.(Errors)
		}
		return this
}
