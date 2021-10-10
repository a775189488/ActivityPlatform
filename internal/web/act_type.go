package web

import (
	code2 "entrytask/internal/common/code"
	"entrytask/internal/common/logger"
	"entrytask/internal/common/perm"
	"entrytask/internal/common/resp"
	"entrytask/internal/service"
	"entrytask/internal/web/metadata"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

type IActTypeWeb interface {
	CreateActType(ctx *gin.Context)
	ListActType(ctx *gin.Context)
	UpdateActType(ctx *gin.Context)
	DeleteActType(ctx *gin.Context)
}

type ActTypeWeb struct {
	Log            logger.ILogger          `inject:""`
	ActTypeService service.IActTypeService `inject:""`
}

func (a *ActTypeWeb) CreateActType(ctx *gin.Context) {
	if !perm.PermCheck(ctx, perm.Admin) {
		a.Log.Errorf("[ActTypeWeb] create activity type fail, Unauthorized")
		resp.RespData(ctx, http.StatusUnauthorized, code2.E5002, "", nil)
		return
	}

	var atReq metadata.ActTypeCreateReq
	if err := ctx.ShouldBind(&atReq); err != nil {
		a.Log.Errorf("[ActTypeWeb] bind atReq object fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6002, "", nil)
		return
	}
	at := atReq.ToActType()
	err := a.ActTypeService.CreateActType(at)
	if err == nil {
		resp.RespSuccess(ctx, "create activity type success", metadata.FormateActTypeCreateResp(at))
	} else {
		handleError(ctx, err)
	}
}

func (a *ActTypeWeb) ListActType(ctx *gin.Context) {
	page, err := com.StrTo(ctx.Query("page")).Int()
	if err != nil {
		a.Log.Errorf("[ActTypeWeb] get page fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	size, err := com.StrTo(ctx.Query("size")).Int()
	if err != nil {
		a.Log.Errorf("[ActTypeWeb] get size fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	data, total, err := a.ActTypeService.ListActType(int32(page), int32(size))
	if err != nil {
		resp.RespData(ctx, http.StatusServiceUnavailable, code2.E9999, "", nil)
		return
	}
	resp.RespSuccess(ctx, "list activity type success", metadata.FormatListActTypeResp(data, total, page, size))

}

func (a *ActTypeWeb) UpdateActType(ctx *gin.Context) {
	if !perm.PermCheck(ctx, perm.Admin) {
		a.Log.Errorf("[ActTypeWeb] update activity type fail, Unauthorized")
		resp.RespData(ctx, http.StatusUnauthorized, code2.E5002, "", nil)
		return
	}

	id, err := com.StrTo(ctx.Param("id")).Int()
	if err != nil {
		a.Log.Errorf("[ActTypeWeb] format activity type fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	var updateReq metadata.ActTypeUpdateReq
	if err := ctx.ShouldBind(&updateReq); err != nil {
		a.Log.Errorf("[ActTypeWeb] bind updateReq object fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	newOne, err := a.ActTypeService.UpdateActType(uint64(id), updateReq.ToActType())
	if err != nil {
		a.Log.Errorf("[ActTypeWeb] update activity type fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	} else {
		resp.RespSuccess(ctx, "update user success", metadata.FormateActTypeUpdateResp(newOne))
	}
}

func (a *ActTypeWeb) DeleteActType(ctx *gin.Context) {
	// todo 不存在也能delete成功
	if !perm.PermCheck(ctx, perm.Admin) {
		a.Log.Errorf("[ActTypeWeb] delete activity type fail, Unauthorized")
		resp.RespData(ctx, http.StatusUnauthorized, code2.E5002, "", nil)
		return
	}
	id, err := com.StrTo(ctx.Param("id")).Int()
	if err != nil {
		a.Log.Errorf("[ActTypeWeb] format activity type id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	err = a.ActTypeService.DeleteActType(uint64(id))
	if err != nil {
		handleError(ctx, err)
	} else {
		resp.RespSuccess(ctx, "delete activity type success", nil)
	}
}
