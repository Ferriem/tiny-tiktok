package model

import (
	"fmt"
	"testing"
	"tiny-tiktok/service/user_service/config"
)

func TestInitDB(t *testing.T) {
	config.InitConfig()
	dns := config.DbDnsInit()
	fmt.Println(dns)
	InitDB()
}
