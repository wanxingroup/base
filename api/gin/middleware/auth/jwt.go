package auth

import (
	"github.com/gin-gonic/gin"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/response"
	log2 "dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/api/gin/util/log"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/auth"
)

// UserId
const UserID = "userId"
const UserInformation = "userInfo"

// 用户token header名称
const RequestUserToken = "Authorization"

/* gin用户jwt认证中间件 */
func UserJwtAuthentication(tokenKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenVal := c.GetHeader(RequestUserToken)
		if len(tokenVal) == 0 {
			tokenVal = c.Query(RequestUserToken)
		}

		userId, info, err := auth.ResolveJWTToken(tokenVal, tokenKey, log2.RequestEntry(c))
		if err != nil {
			c.Abort()

			if auth.IsTokenExpired(err) {
				response.TokenExpired(c)
				return
			}

			response.TokenInvalid(c)
			return
		}

		c.Set(UserID, userId)
		c.Set(UserInformation, info)
		c.Next()
	}
}
