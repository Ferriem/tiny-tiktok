package router

import (
	"tiny-tiktok/api_router/internal/handler"
	"tiny-tiktok/api_router/router/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(serveInstance map[string]interface{}) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.ServeMiddleware(serveInstance), middleware.ErrorMiddleware())
	//test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	baseGroup := r.Group("/tiktok")
	{
		//user
		baseGroup.POST("/user/register/", handler.UserRegister)
		baseGroup.POST("/user/login/", handler.UserLogin)
		baseGroup.GET("/user/", middleware.JWTMiddleware(), handler.UserInfo)

		//social
		relationGroup := baseGroup.Group("/relation")
		relationGroup.Use(middleware.JWTMiddleware())
		{
			relationGroup.POST("/action/", handler.FollowAction)
			relationGroup.GET("/follow/list/", handler.GetFollowList)
			relationGroup.GET("/follower/list/", handler.GetFollowerList)
			relationGroup.GET("/friend/list/", handler.GetFriendList)
		}
		messageGroup := baseGroup.Group("/message")
		messageGroup.Use(middleware.JWTMiddleware())
		{
			messageGroup.POST("/send/", handler.PostMessage)
			messageGroup.GET("/chat/", handler.GetMessage)
		}
	}
	return r
}
