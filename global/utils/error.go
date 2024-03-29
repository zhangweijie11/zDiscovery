package utils

import "fmt"

type Error struct {
	code    int      `json:"code"`
	msg     string   `json:"msg"`
	details []string `json:"details"`
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("code %d exists!", code))
	}
	codes[code] = msg
	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, error info: %s", e.code, e.msg)
}

func (e *Error) Code() int {
	return e.code
}

var (
	Success     = NewError(200, "success")
	NotModified = NewError(304, "app not modified")
	ParamError  = NewError(400, "request param error")
	NotFound    = NewError(404, "not found")
	Conflict    = NewError(409, "conflict")
	ServerError = NewError(500, "service internal error")
)
