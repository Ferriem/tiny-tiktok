package model

import (
	"fmt"
	"testing"
	"tiny-tiktok/service/video_service/config"
)

func TestInitDb(t *testing.T) {
	config.InitConfig()
	dns := config.DbDnsInit()
	fmt.Println(dns)
	InitDB()
}
