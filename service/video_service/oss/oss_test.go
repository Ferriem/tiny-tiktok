package oss

import (
	"fmt"
	"os"
	"testing"
	"tiny-tiktok/api_router/pkg/logger"
	"tiny-tiktok/service/video_service/config"
)

func TestUpload(t *testing.T) {
	config.InitConfig()
	filePath := "/Users/ferriem/Downloads/若能绽放光芒.mp4"

	fileByte, _ := os.ReadFile(filePath)
	logger.Log.Info("fileByte: ", fileByte)
	err := UpLoad("若能绽放光芒.mp4", fileByte)
	fmt.Println(err)
}

func TestDelete(t *testing.T) {
	config.InitConfig()
	err := Delete("1.png")
	fmt.Println(err)
}
