package web

import (
	"entrytask/internal/common"
	code2 "entrytask/internal/common/code"
	"entrytask/internal/common/logger"
	"entrytask/internal/common/perm"
	"entrytask/internal/common/resp"
	"entrytask/internal/common/utils"
	"entrytask/internal/conf"
	"entrytask/internal/service"
	"entrytask/internal/web/metadata"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/unknwon/com"
)

type IUserWeb interface {
	Test(ctx *gin.Context)
	Register(ctx *gin.Context)
	ListUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

type UserWeb struct {
	Log         logger.ILogger       `inject:""`
	UserService service.IUserSerivce `inject:""`
}

func (u *UserWeb) Test(ctx *gin.Context) {
	resp.RespSuccess(ctx, "done!", nil)
}

func (u *UserWeb) DeleteUser(ctx *gin.Context) {
	// todo 做权限检查，只允许管理员删除
	id, err := com.StrTo(ctx.Param("id")).Int()
	if err != nil {
		u.Log.Errorf("[UserWeb] get user id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	err = u.UserService.DeleteUser(uint64(id))
	if err != nil {
		handleError(ctx, err)
	} else {
		resp.RespSuccess(ctx, "delete user success", nil)
	}
}

func (u *UserWeb) UpdateUser(ctx *gin.Context) {
	// todo 做权限检查，只允许用户自己的账户
	id, err := com.StrTo(ctx.Param("id")).Int()
	if err != nil {
		u.Log.Errorf("[UserWeb] get user id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	var updateReq metadata.UpdateUserReq
	if err := ctx.ShouldBind(&updateReq); err != nil {
		u.Log.Errorf("[UserWeb] bind updateReq object fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	newOne, err := u.UserService.UpdateUser(uint64(id), updateReq.ToUser())
	if err != nil {
		u.Log.Errorf("[UserWeb] get user id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	} else {
		resp.RespSuccess(ctx, "update user success", metadata.FormatGetUserResp(newOne))
	}
}

func (u *UserWeb) GetUser(ctx *gin.Context) {
	id, err := com.StrTo(ctx.Param("id")).Int()
	if err != nil {
		u.Log.Errorf("[UserWeb] get user id fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	result, err := u.UserService.GetUserById(uint64(id))
	if err != nil {
		handleError(ctx, err)
	} else {
		resp.RespSuccess(ctx, "get user success", metadata.FormatGetUserResp(result))
	}
}

func (u *UserWeb) Register(ctx *gin.Context) {
	var userReq metadata.UserRegisterReq
	// todo 参数校验
	if err := ctx.ShouldBind(&userReq); err != nil {
		u.Log.Errorf("[UserWeb] bind userReq object fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	user := userReq.ToUser()
	err := u.UserService.CreateUser(user)
	if err == nil {
		resp.RespSuccess(ctx, "register success", metadata.FormatUserRegisterResp(user))
	} else {
		handleError(ctx, err)
	}
}

func (u *UserWeb) ListUser(ctx *gin.Context) {
	if !perm.PermCheck(ctx, perm.Admin) {
		u.Log.Errorf("[UserWeb] List user fail, Unauthorized")
		resp.RespData(ctx, http.StatusUnauthorized, code2.E5002, "", nil)
		return
	}

	// todo 现在成为必传了
	page, err := com.StrTo(ctx.Query("page")).Int()
	if err != nil {
		u.Log.Errorf("[UserWeb] get page fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	size, err := com.StrTo(ctx.Query("size")).Int()
	if err != nil {
		u.Log.Errorf("[UserWeb] get size fail, err: %v", err)
		resp.RespData(ctx, http.StatusBadRequest, code2.E6001, "", nil)
	}
	users, total, err := u.UserService.ListUser(int32(page), int32(size))
	if err != nil {
		resp.RespData(ctx, http.StatusServiceUnavailable, code2.E9999, "", nil)
		return
	}
	resp.RespSuccess(ctx, "list user success", metadata.FormatListUserResp(users, total, page, size))
}

func (u *UserWeb) UploadImage(ctx *gin.Context) {
	u.Log.Infof("[UserWeb]Upload image")
	file, header, err := ctx.Request.FormFile("profile")
	if err != nil {
		resp.RespData(ctx, http.StatusServiceUnavailable, code2.E9999, "", nil)
	}
	filename := header.Filename
	// todo 是否需要检查文件后缀名和文件大小增加安全性
	suffix := filename[strings.LastIndex(filename, "."):]

	sum := utils.Md5UploadFile(file)
	pathAppend := common.HeadPicturePath + strconv.Itoa(time.Now().Year()) + "_" + sum + suffix
	var path = conf.Config.App.FilePath + pathAppend
	out, err := os.Create(path)
	if err != nil {
		u.Log.Errorf("create file(%s) fail, err: %s", path, err)
		resp.RespData(ctx, http.StatusServiceUnavailable, code2.E9999, "", nil)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			u.Log.Errorf("upload file close err:%v", err)
		}
	}(file)

	_, errSeek := file.Seek(0, 0) // 重置文件指针
	if errSeek != nil {
		u.Log.Errorf("seek file err:%v", errSeek)
		resp.RespData(ctx, http.StatusServiceUnavailable, code2.E9999, "", nil)
		return
	}
	_, err = io.Copy(out, file)
	if err != nil {
		if os.IsExist(err) {
			u.Log.Errorf("file exist, file name:%s", out.Name())
		} else {
			u.Log.Errorf("storage file err:%v", err)
		}
		resp.RespData(ctx, http.StatusServiceUnavailable, code2.E9999, "", nil)
		return
	}
	u.Log.Infof("upload file success, filename:%s, pathAppend:%s, url:%s", filename, pathAppend, path)

	resp.RespSuccess(ctx, "upload success", &metadata.UploadHeadpicResp{FilePath: pathAppend})
}
