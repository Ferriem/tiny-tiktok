package main

import (
	"tiny-tiktok/service/social_service/config"
	"tiny-tiktok/service/social_service/discovery"
	"tiny-tiktok/service/social_service/internal/model"
	"tiny-tiktok/service/social_service/pkg/cache"
)

func main() {
	config.InitConfig()
	model.InitDB()
	cache.InitRedis()
	go cache.TimeSync()
	discovery.AutoRegister()
}
