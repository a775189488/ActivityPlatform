package resp

import (
	"net/http"
	"strconv"

	"entrytask/internal/common/code"

	"entrytask/internal/conf"
	"github.com/gin-gonic/gin"
)

//ResponseData 数据返回结构体
type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespSuccess(c *gin.Context, msg string, data interface{}) {
	resp := ResponseData{
		Code:    int(code.E0000),
		Message: msg,
		Data:    data,
	}
	RespJSON(c, http.StatusOK, resp)
}

func RespData(c *gin.Context, httpCode int, errCode code.ErrCode, msg string, data interface{}) {
	if msg == "" {
		msg = code.GetErrMsg(errCode)
	}
	resp := ResponseData{
		Code:    int(errCode),
		Message: msg,
		Data:    data,
	}
	RespJSON(c, httpCode, resp)
}

func RespJSON(c *gin.Context, httpCode int, resp interface{}) {
	if resp == nil {
		resp = map[string]interface{}{}
	}
	c.JSON(httpCode, resp)
	c.Abort()
}

func GetPage(c *gin.Context) (page, pagesize int) {
	page, _ = strconv.Atoi(c.Query("page"))
	pagesize, _ = strconv.Atoi(c.Query("limit"))
	if pagesize == 0 {
		pagesize = conf.Config.App.PageSize
	}
	if page == 0 {
		page = 1
	}
	return
}
