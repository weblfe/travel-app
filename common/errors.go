package common

import "fmt"

type Errors interface {
		Code() int
		Msg() string
		Parent() Errors
		error
}

type ErrorImpl struct {
		code   int
		msg    string
		parent Errors
}

// code string
// msg  int
func NewErrors(args ...interface{}) Errors {
		var err = new(ErrorImpl)
		err.code = -1
		err.init(args...)
		return err
}

// err Errors
// code string
// msg  int
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
				}
				if msg, ok := arg.(string); ok && this.msg == "" {
						this.msg = msg
				}
				if code, ok := arg.(int); ok && this.code == -1 {
						this.code = code
				}
		}
}

func (this *ErrorImpl) Error() string {
		return fmt.Sprintf(`{"code":%d,"msg":"%s"}`, this.Code(), this.Msg())
}

func (this *ErrorImpl) Parent() Errors {
		return this.parent
}

func (this *ErrorImpl) Code() int {
		return this.code
}

func (this *ErrorImpl) Msg() string {
		return this.msg
}
