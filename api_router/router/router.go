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
		baseGroup.GET("/feed/", handler.Feed)

		//user
		baseGroup.POST("/user/register/", handler.UserRegister)
		baseGroup.POST("/user/login/", handler.UserLogin)
		baseGroup.GET("/user/", middleware.JWTMiddleware(), handler.UserInfo)

		//video
		publishGroup := baseGroup.Group("/publish")
		publishGroup.Use(middleware.JWTMiddleware())
		{
			publishGroup.POST("/action/", handler.PublishAction)
			publishGroup.GET("/list/", handler.PublishList)
		}
		favoriteGroup := baseGroup.Group("/favorite")
		favoriteGroup.Use(middleware.JWTMiddleware())
		{
			favoriteGroup.POST("/action/", handler.FavoriteAction)
			favoriteGroup.GET("/list/", handler.FavoriteList)
		}
		commentGroup := baseGroup.Group("/comment")
		commentGroup.Use(middleware.JWTMiddleware())
		{
			commentGroup.POST("/action/", handler.CommentAction)
			commentGroup.GET("/list/", handler.CommentList)
		}

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
