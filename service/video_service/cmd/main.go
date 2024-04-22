package main

import (
	"tiny-tiktok/service/video_service/config"
	"tiny-tiktok/service/video_service/discovery"
	"tiny-tiktok/service/video_service/internal/handler"
	"tiny-tiktok/service/video_service/internal/model"
	"tiny-tiktok/service/video_service/pkg/cache"
)

func main() {
	config.InitConfig()
	model.InitDB()
	cache.InitRedis()
	go func() {
		handler.PublishVideo()
	}()
	discovery.AutoRegister()
}
