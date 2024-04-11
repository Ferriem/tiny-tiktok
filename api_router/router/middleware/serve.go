package middleware

import "github.com/gin-gonic/gin"

func ServeMiddleware(serveInstance map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Keys = make(map[string]interface{})
		for key, value := range serveInstance {
			c.Keys[key] = value
		}
		c.Next()
	}
}
