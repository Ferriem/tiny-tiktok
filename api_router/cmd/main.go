package main

import (
	"fmt"
	"net/http"
	"time"
	"tiny-tiktok/api_router/config"
	"tiny-tiktok/api_router/discovery"
	"tiny-tiktok/api_router/pkg/logger"
	"tiny-tiktok/api_router/router"

	"github.com/spf13/viper"
)

func main() {

	config.InitConfig()
	resolver := discovery.Resolver()

	r := router.InitRouter(resolver)

	server := &http.Server{
		Addr:           viper.GetString("server.address"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Server is running on", viper.GetString("server.address"))
	err := server.ListenAndServe()
	if err != nil {
		logger.Log.Fatal("Start failed...")
	}
}
