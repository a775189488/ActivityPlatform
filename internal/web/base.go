package web

import (
	code2 "entrytask/internal/common/code"
	error2 "entrytask/internal/common/error"
	"entrytask/internal/common/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func handleError(ctx *gin.Context, err error) {
	mError, ok := err.(error2.MyError)
	if !ok {
		resp.RespData(ctx, http.StatusServiceUnavailable, code2.E9999, "", nil)
	} else {
		resp.RespData(ctx, http.StatusBadRequest, mError.Code, "", nil)
	}
}
