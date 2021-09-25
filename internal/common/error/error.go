package error

import (
	code2 "entrytask/internal/common/code"
	"errors"
)

type MyError struct {
	error
	Code code2.ErrCode
}

func NewMyError(code code2.ErrCode, msg string) MyError {
	if msg == "" {
		msg = code2.GetErrMsg(code)
	}
	return MyError{
		errors.New(msg),
		code,
	}
}

var UserNameConflictError = NewMyError(code2.E1001, "")
var UserNotFoundError = NewMyError(code2.E1002, "")
var UserPasswordError = NewMyError(code2.E1003, "")

var ActTypeNotFoundError = NewMyError(code2.E3001, "")
var ActTypeParentNotFoundError = NewMyError(code2.E3002, "")
