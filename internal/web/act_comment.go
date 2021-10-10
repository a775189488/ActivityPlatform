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

type IActCommentWeb interface {
	CreateActComment(ctx *gin.Context)
	UpdateActComment(ctx *gin.Context)
	ListActCommentByActId(ctx *gin.Context)
	DeleteActComment(ctx *gin.Context)
}

type ActCommentWeb struct {
	Log               logger.ILogger             `inject:""`
	ActCommentService service.IActCommentService `inject:""`
}

func (a *ActCommentWeb) CreateActComment(ctx *gin.Context) {
	var acReq metadata.ActCommentCreateReq
	if err := ctx.ShouldBind(&acReq); err != nil {
		a.Log.Errorf("[ActCommentWeb] bind acReq object fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6002, "", nil)
		return
	}
	ac := acReq.ToActComment(perm.GetUserIdFromCtx(ctx))
	if err := a.ActCommentService.CreateActComment(ac); err != nil {
		handleError(ctx, err)
	} else {
		resp.RespSuccess(ctx, "create activity comment success", metadata.FormatActCommentCreateResp(ac))
	}
}

func (a *ActCommentWeb) UpdateActComment(ctx *gin.Context) {
	// todo 只能改自己的评论
	id, err := com.StrTo(ctx.Param("id")).Int()
	if err != nil {
		a.Log.Errorf("[ActCommentWeb] format activity comment id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	var updateReq metadata.ActCommentUpdateReq
	if err := ctx.ShouldBind(&updateReq); err != nil {
		a.Log.Errorf("[ActCommentWeb] bind updateReq object fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	newOne, err := a.ActCommentService.UpdateActComment(uint64(id), updateReq.ToActComment())
	if err != nil {
		a.Log.Errorf("[ActCommentWeb] update activity comment fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	} else {
		resp.RespSuccess(ctx, "update user success", metadata.FormatActCommentUpdateResp(newOne))
	}
}

func (a *ActCommentWeb) ListActCommentByActId(ctx *gin.Context) {
	id, err := com.StrTo(ctx.Param("actid")).Int()
	if err != nil {
		a.Log.Errorf("[ActCommentWeb] format activity comment id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	page, err := com.StrTo(ctx.Query("page")).Int()
	if err != nil {
		a.Log.Errorf("[ActCommentWeb] get page fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	size, err := com.StrTo(ctx.Query("size")).Int()
	if err != nil {
		a.Log.Errorf("[ActCommentWeb] get size fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	data, total, err := a.ActCommentService.ListActCommentByActId(int32(page), int32(size), uint64(id))
	if err != nil {
		handleError(ctx, err)
	} else {
		resp.RespSuccess(ctx, "list activity comment success", metadata.FormatActCommentListResp(data, total, page, size))
	}
}

func (a *ActCommentWeb) DeleteActComment(ctx *gin.Context) {
	if !perm.PermCheck(ctx, perm.Admin) {
		a.Log.Errorf("[ActCommentWeb] delete activity comment fail, Unauthorized")
		resp.RespData(ctx, http.StatusUnauthorized, code2.E5002, "", nil)
		return
	}
	id, err := com.StrTo(ctx.Param("id")).Int()
	if err != nil {
		a.Log.Errorf("[ActCommentWeb] format activity comment id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	err = a.ActCommentService.DeleteActComment(uint64(id))
	if err != nil {
		handleError(ctx, err)
	} else {
		resp.RespSuccess(ctx, "delete activity comment success", nil)
	}
}
