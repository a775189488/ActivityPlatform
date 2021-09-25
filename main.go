package main

import (
	"fmt"
	_ "github.com/appleboy/gin-jwt/v2"
	_ "github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm"
	_ "go.uber.org/zap"
	_ "go.uber.org/zap/zapcore"
	_ "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	fmt.Println(r)
}
