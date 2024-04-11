package oss

import (
	"fmt"
	"os"
	"testing"
	"tiny-tiktok/service/video_service/config"
)

func TestUpload(t *testing.T) {
	config.InitConfig()
	filePath := "/Users/ferriem/Downloads/91953835_p1.png"

	fileByte, _ := os.ReadFile(filePath)
	err := UpLoad("1.png", fileByte)
	fmt.Println(err)
}

func TestDelete(t *testing.T) {
	config.InitConfig()
	err := Delete("1.png")
	fmt.Println(err)
}
