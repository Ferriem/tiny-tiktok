package config

import (
	"fmt"
	"path"
	"runtime"
	"testing"

	"github.com/spf13/viper"
)

func TestDbDnsInit(t *testing.T) {
	InitConfig()
	fmt.Printf("host is : %v \n", viper.GetString("mysql.host"))
	dns := DbDnsInit()
	fmt.Printf("dns is : %v", dns)
}

func TestInitConfig(t *testing.T) {
	_, filePath, _, _ := runtime.Caller(0)

	currentDir := path.Dir(filePath)

	fmt.Println(currentDir)
}
