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
var UserJoinDupActivityError = NewMyError(code2.E1004, "")
var UserLeaveNotJoinActivityError = NewMyError(code2.E1005, "")

var ActCreateTypeNotFoundError = NewMyError(code2.E2001, "")

var ActTypeNotFoundError = NewMyError(code2.E3001, "")
var ActTypeParentNotFoundError = NewMyError(code2.E3002, "")
var ActTypeDeleteError = NewMyError(code2.E3003, "")

var ActCommentCreateActNotExistError = NewMyError(code2.E4001, "")
var ActCommentCreateUserNotExistError = NewMyError(code2.E4002, "")
