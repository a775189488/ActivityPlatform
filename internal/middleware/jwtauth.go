package middleware

import (
	"entrytask/internal/common/logger"
	"net/http"
	"time"

	"encoding/base64"
	code2 "entrytask/internal/common/code"
	"entrytask/internal/common/resp"
	"entrytask/internal/common/utils"
	"entrytask/internal/conf"
	"entrytask/internal/model/vo"
	"entrytask/internal/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type UserLoginReq struct {
	Username string `json:"username";validate:"min=6,max=10"`
	Password string `json:"password";validate:"min=6,max=10"`
}

type Jwt struct {
	Log         logger.ILogger       `inject:""`
	UserService service.IUserSerivce `inject:""`
}

type JwtAuthorizator func(data interface{}, c *gin.Context) bool

func (j *Jwt) GinJWTMiddlewareInit(jwtAuthorizator JwtAuthorizator) (authMiddleware *jwt.GinJWTMiddleware) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Minute * 60,
		MaxRefresh:  time.Hour,
		IdentityKey: conf.Config.App.IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*vo.UserLoginVo); ok {
				//maps the claims in the JWT
				return jwt.MapClaims{
					"id":        v.Id,
					"username":  v.Username,
					"role":      v.Role,
					"aliasname": v.Aliasname,
					"email":     v.Email,
					"headpic":   v.Headpic,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			userInfo := jwt.ExtractClaims(c)
			userId := uint64(userInfo["id"].(float64))
			username := userInfo["username"].(string)
			role := int(userInfo["id"].(float64))
			aliasname := userInfo["aliasname"].(string)
			email := userInfo["email"].(string)
			headpic := userInfo["headpic"].(string)
			//Set the identity
			return &vo.UserLoginVo{
				Id:        userId,
				Username:  username,
				Role:      role,
				Aliasname: aliasname,
				Email:     email,
				Headpic:   headpic,
			}
		},
		// 登录逻辑
		Authenticator: func(c *gin.Context) (interface{}, error) {
			//handles the login logic. On success LoginResponse is called, on failure Unauthorized is called
			var user UserLoginReq
			if err := c.ShouldBind(&user); err != nil {
				j.Log.Errorf("[jwt] bind userlogin object fail, err: %v", err)
				return "", jwt.ErrMissingLoginValues
			}
			username := user.Username
			password := user.Password
			// 对password解密
			pwd, _ := base64.StdEncoding.DecodeString(password)
			realPwd := utils.AesCbcDecrypt(pwd, []byte(conf.Config.App.AesKey))
			// md5加密
			if err := j.UserService.CheckUser(username, utils.Md5Encrypt(realPwd)); err != nil {
				j.Log.Errorf("[jwt] check user(%s) fail, err: %v", username, err)
				return nil, jwt.ErrFailedAuthentication
			}
			userInfo, findUserErr := j.UserService.GetUserInfoByLoginName(username)
			if findUserErr != nil {
				j.Log.Errorf("[jwt] get user by LoginName(%s) fail, err: %v", username, findUserErr)
				return nil, jwt.ErrFailedAuthentication
			}
			return vo.FormatUserLoginVo(userInfo), nil
		},
		//receives identity and handles authorization logic
		Authorizator: jwtAuthorizator,
		//handles unauthorized logic
		Unauthorized: func(c *gin.Context, code int, message string) {
			resp.RespData(c, http.StatusUnauthorized, code2.E5001, message, nil)
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			data := map[string]interface{}{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			}
			resp.RespSuccess(c, "login success", data)
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			data := map[string]interface{}{
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			}
			resp.RespSuccess(c, "refresh success", data)
		},
		LogoutResponse: func(c *gin.Context, code int) {
			resp.RespSuccess(c, "logout success", nil)
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		j.Log.DPanicf("JWT Error:" + err.Error())
	}
	return
}

func Authorizator(data interface{}, c *gin.Context) bool {
	return true
}

//NoRouteHandler 404 handler
func NoRouteHandler(c *gin.Context) {
	resp.RespData(c, http.StatusNotFound, code2.E6001, "", nil)
}
