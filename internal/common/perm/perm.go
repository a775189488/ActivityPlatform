package perm

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type PermType int

const (
	Default = 0
	Admin   = 1
)

func PermCheck(ctx *gin.Context, needPerm PermType) bool {
	userInfo := jwt.ExtractClaims(ctx)
	role := int(userInfo["role"].(float64))
	if PermType(role) != needPerm {
		return false
	}
	return true
}

func GetUserIdFromCtx(ctx *gin.Context) uint64 {
	userInfo := jwt.ExtractClaims(ctx)
	if userInfo["id"] == nil {
		return 0
	}
	userId := uint64(int(userInfo["id"].(float64)))
	return userId
}
