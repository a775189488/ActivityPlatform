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
