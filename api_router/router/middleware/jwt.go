package middleware

import (
	"net/http"
	"time"
	"tiny-tiktok/api_router/pkg/auth"
	"tiny-tiktok/utils/exceptions"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = exceptions.SUCCESS
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}

		claims, err := auth.ParseToken(token)
		if err != nil {
			code = exceptions.UnAuth
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = exceptions.TokenTimeout
		}

		if code != exceptions.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"StatusCode": code,
				"Msg":        exceptions.GetMsg(code),
			})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
