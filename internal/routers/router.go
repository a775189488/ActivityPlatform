package routers

import (
	db2 "entrytask/internal/common/db"
	"entrytask/internal/common/logger"
	"entrytask/internal/common/resp"
	"entrytask/internal/conf"
	"entrytask/internal/middleware"
	"entrytask/internal/repository"
	"entrytask/internal/service"
	"entrytask/internal/web"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"log"
)

func InitRouter() *gin.Engine {
	if conf.Config.App.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	Configure(r)
	return r
}

func Configure(r *gin.Engine) {
	var userApi web.UserWeb
	var myjwt middleware.Jwt
	db := db2.DbMysql{}
	zap := logger.Logger{}
	var injector inject.Graph
	if err := injector.Provide(
		&inject.Object{Value: &userApi},
		&inject.Object{Value: &db},
		&inject.Object{Value: &zap},
		&inject.Object{Value: &myjwt},
		&inject.Object{Value: &repository.UserRepo{}},
		&inject.Object{Value: &service.UserService{}},
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
	r.NoRoute(authMiddleware.MiddlewareFunc(), middleware.NoRouteHandler)
	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/refresh_token", authMiddleware.RefreshHandler)
	r.GET("/ping", ping)
	r.POST("/register", userApi.Register)

	// todo 使用 jwt 无法注销token
	r.POST("/logout", authMiddleware.LogoutHandler).Use(authMiddleware.MiddlewareFunc())
	v1 := r.Group("/v1")
	v1.Use(authMiddleware.MiddlewareFunc())
	{
		v1.GET("/user", userApi.ListUser)
		v1.GET("/user/:id", userApi.GetUser)
		v1.PUT("/user/:id", userApi.UpdateUser)
		v1.POST("/upload", userApi.UploadImage)
		v1.DELETE("/user/:id", userApi.DeleteUser)
	}
}

func ping(ctx *gin.Context) {
	resp.RespSuccess(ctx, "Pong!", nil)
}
