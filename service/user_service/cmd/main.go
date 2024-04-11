package main

import (
	"tiny-tiktok/service/user_service/config"
	"tiny-tiktok/service/user_service/discovery"
	"tiny-tiktok/service/user_service/internal/model"
	"tiny-tiktok/service/user_service/pkg/cache"
)

func main() {
	config.InitConfig()
	cache.InitRedis()
	model.InitDB()
	discovery.AutoRegister()
}
