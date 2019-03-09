package exception

import "fmt"

type Exception interface {
	GetCode() int
	GetMsg() string
}

type exception struct {
	code int
	msg  string
}

func (e *exception) GetCode() int { return e.code }

func (e *exception) GetMsg() string { return e.msg }

func (e *exception) FormatErrStr() string {
	return fmt.Sprintf("[%d] %v", e.code, e.msg)
}

func New(code int, msg string) Exception {
	return &exception{code, msg}
}

func NewFromErr(code int, err error) Exception {
	return &exception{code: code, msg: err.Error()}
}
