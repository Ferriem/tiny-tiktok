package middleware

import (
	"fmt"
	"net/http"
	"tiny-tiktok/utils/exceptions"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			r := recover()
			if r != nil {
				c.JSON(http.StatusOK, gin.H{
					"status_code": exceptions.ERROR,
					"status_msg":  fmt.Sprintf("%s", r),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
