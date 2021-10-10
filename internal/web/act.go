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

type IActivityWeb interface {
	CreateActivity(ctx *gin.Context)
	ListActivity(ctx *gin.Context)
	UpdateActivity(ctx *gin.Context)
	DeleteActivity(ctx *gin.Context)
	GetActivityDetail(ctx *gin.Context)
	GetActivityUser(ctx *gin.Context)
}

type ActivityWeb struct {
	Log             logger.ILogger      `inject:""`
	ActivityService service.IActService `inject:""`
}

func (a *ActivityWeb) GetActivityDetail(ctx *gin.Context) {
	id, err := com.StrTo(ctx.Param("actid")).Int()
	if err != nil {
		a.Log.Errorf("[ActWeb] get activity id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
		return
	}

	result, err := a.ActivityService.GetActivityDetail(uint64(id))
	if err != nil {
		handleError(ctx, err)
		return
	}

	isJoin := false
	userId := perm.GetUserIdFromCtx(ctx)
	if userId != 0 {
		isJoin, err = a.ActivityService.CheckUserJoinActivity(uint64(id), perm.GetUserIdFromCtx(ctx))
		if err != nil {
			handleError(ctx, err)
			return
		}
	}

	resp.RespSuccess(ctx, "get activity success", metadata.FormatGetActDetailResp(result, isJoin))
}

func (a *ActivityWeb) CreateActivity(ctx *gin.Context) {
	var actReq metadata.ActCreateReq
	if err := ctx.ShouldBind(&actReq); err != nil {
		a.Log.Errorf("[ActWeb] bind actReq object fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6002, "", nil)
		return
	}
	act := actReq.ToActivity(perm.GetUserIdFromCtx(ctx))
	if err := a.ActivityService.CreateActivity(act); err != nil {
		handleError(ctx, err)
	} else {
		resp.RespSuccess(ctx, "create activity success", metadata.FormatActCreateResp(act))
	}
}

func (a *ActivityWeb) ListActivity(ctx *gin.Context) {
	var actReq metadata.ListActReq
	if err := ctx.BindQuery(&actReq); err != nil {
		a.Log.Errorf("[ActWeb] bind actReq object fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6002, "", nil)
		return
	}
	acts, err := a.ActivityService.ListActivity(actReq.Page, actReq.Size, actReq.ToWhere())
	if err != nil {
		handleError(ctx, err)
		return
	}

	userId := perm.GetUserIdFromCtx(ctx)
	joinActs := make([]uint64, 0)
	if userId != 0 {
		actIds := make([]uint64, 0)
		for _, a := range acts {
			actIds = append(actIds, a.Id)
		}
		joinActs, err = a.ActivityService.CheckUserJoinActivities(actIds, perm.GetUserIdFromCtx(ctx))
		if err != nil {
			handleError(ctx, err)
			return
		}
	}
	resp.RespSuccess(ctx, "list activity success", metadata.FormatListActResp(acts, int(actReq.Page), int(actReq.Size), joinActs))
}

func (a *ActivityWeb) UpdateActivity(ctx *gin.Context) {
	id, err := com.StrTo(ctx.Param("id")).Int()
	if err != nil {
		a.Log.Errorf("[ActWeb] format activity id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	// todo 只允许修改自己创建的活动，或者管理员修改
	var updateReq metadata.ActUpdateReq
	if err := ctx.ShouldBind(&updateReq); err != nil {
		a.Log.Errorf("[ActWeb] bind updateReq object fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	newOne, err := a.ActivityService.UpdateActivity(uint64(id), updateReq.ToActivity())
	if err != nil {
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	} else {
		resp.RespSuccess(ctx, "update activity success", metadata.FormatActUpdateResp(newOne))
	}
}

func (a *ActivityWeb) DeleteActivity(ctx *gin.Context) {
	// todo 只允许删除自己创建的活动，或者管理员删除
	id, err := com.StrTo(ctx.Param("id")).Int()
	if err != nil {
		a.Log.Errorf("[ActWeb] format activity id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
		return
	}
	err = a.ActivityService.DeleteActivity(uint64(id))
	if err != nil {
		handleError(ctx, err)
	} else {
		resp.RespSuccess(ctx, "delete activity success", nil)
	}
}

func (a *ActivityWeb) GetActivityUser(ctx *gin.Context) {
	id, err := com.StrTo(ctx.Param("actid")).Int()
	if err != nil {
		a.Log.Errorf("[ActWeb] format activity id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
		return
	}
	var actUserReq metadata.ListActUserReq
	if err := ctx.BindQuery(&actUserReq); err != nil {
		a.Log.Errorf("[ActWeb] bind actUserReq object fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6002, "", nil)
		return
	}
	users, err := a.ActivityService.GetActivityUser(actUserReq.Page, actUserReq.Size, uint64(id))
	if err != nil {
		handleError(ctx, err)
	} else {
		resp.RespSuccess(ctx, "get activity users success", metadata.FormateGetActivityUser(users, int(actUserReq.Page), int(actUserReq.Size)))
	}
}
