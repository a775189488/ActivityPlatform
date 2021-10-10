package routers

import (
	"entrytask/internal/common"
	db2 "entrytask/internal/common/db"
	"entrytask/internal/common/logger"
	"entrytask/internal/common/resp"
	"entrytask/internal/conf"
	"entrytask/internal/middleware"
	"entrytask/internal/repository"
	"entrytask/internal/service"
	"entrytask/internal/web"
	"github.com/facebookgo/inject"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func InitRouter() *gin.Engine {

	r := gin.New()

	if conf.Config.App.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		// 如果是测试版，那么打开pprof用于性能问题排查
		pprof.Register(r)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	Configure(r)

	return r
}

func Configure(r *gin.Engine) {
	var userApi web.UserWeb
	var actTypeApi web.ActTypeWeb
	var actCommentApi web.ActCommentWeb
	var actApi web.ActivityWeb
	var myjwt middleware.Jwt
	db := db2.DbMysql{}
	zap := logger.Logger{}
	var injector inject.Graph
	if err := injector.Provide(
		&inject.Object{Value: &actApi},
		&inject.Object{Value: &userApi},
		&inject.Object{Value: &actTypeApi},
		&inject.Object{Value: &actCommentApi},
		&inject.Object{Value: &db},
		&inject.Object{Value: &zap},
		&inject.Object{Value: &myjwt},
		&inject.Object{Value: &repository.UserRepo{}},
		&inject.Object{Value: &repository.ActivityRepo{}},
		&inject.Object{Value: &repository.ActivityTypeRepo{}},
		&inject.Object{Value: &repository.ActivityUserRepo{}},
		&inject.Object{Value: &repository.ActivityCommentRepo{}},
		&inject.Object{Value: &service.ActTypeService{}},
		&inject.Object{Value: &service.ActService{}},
		&inject.Object{Value: &service.UserService{}},
		&inject.Object{Value: &service.ActCommentService{}},
	); err != nil {
		log.Fatal("inject fatal: ", err)
	}
	if err := injector.Populate(); err != nil {
		log.Fatal("injector fatal: ", err)
	}

	zap.Init()
	if err := db.Connect(); err != nil {
		log.Fatal("db fatal:", err)
	}
	var authMiddleware = myjwt.GinJWTMiddlewareInit(middleware.Authorizator)

	r.StaticFS("/file", http.Dir(conf.Config.App.FilePath+common.HeadPicturePath))
	r.NoRoute(authMiddleware.MiddlewareFunc(), middleware.NoRouteHandler)
	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/refresh_token", authMiddleware.RefreshHandler)
	r.GET("/ping", ping)
	r.POST("/register", userApi.Register)

	r.GET("/act", actApi.ListActivity)
	r.GET("/act/:actid", actApi.GetActivityDetail)

	r.GET("/actType", actTypeApi.ListActType)
	r.GET("/act/:actid/comment", actCommentApi.ListActCommentByActId)
	r.GET("/act/:actid/user", actApi.GetActivityUser)

	// todo 使用 jwt 无法注销token
	r.POST("/logout", authMiddleware.LogoutHandler).Use(authMiddleware.MiddlewareFunc())
	v1 := r.Group("/v1", authMiddleware.MiddlewareFunc())
	{
		// user
		v1.GET("/user", userApi.ListUser)
		v1.GET("/user/:id", userApi.GetUser)
		v1.PUT("/user/:id", userApi.UpdateUser)
		v1.POST("/upload", userApi.UploadImage)
		v1.DELETE("/user/:id", userApi.DeleteUser)
		v1.POST("/user/activity/:id", userApi.JoinActivity)
		v1.DELETE("/user/activity/:id", userApi.LeaveActivity)
		v1.GET("/user/:id/activity", userApi.GetUserActivity)

		//activity type
		v1.PUT("/actType/:id", actTypeApi.UpdateActType)
		v1.DELETE("/actType/:id", actTypeApi.DeleteActType)
		v1.POST("/actType", actTypeApi.CreateActType)

		//activity comment
		v1.PUT("/comment/:id", actCommentApi.UpdateActComment)
		v1.POST("/comment", actCommentApi.CreateActComment)
		v1.DELETE("/comment/:id", actCommentApi.DeleteActComment)

		v1.POST("/act", actApi.CreateActivity)
		v1.GET("/act/:actid", actApi.GetActivityDetail)
		v1.GET("/act", actApi.ListActivity)
		v1.PUT("/act/:id", actApi.UpdateActivity)
		v1.DELETE("/act/:id", actApi.DeleteActivity)
	}
}

func ping(ctx *gin.Context) {
	resp.RespSuccess(ctx, "Pong!", nil)
}
