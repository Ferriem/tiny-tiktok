package config

import (
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Mysql struct {
		Host     string
		Port     string
		Username string
		Password string
		Database string
	}
}

var ConfigData *Config

func InitConfig() {
	_, filePath, _, _ := runtime.Caller(0)

	currenDir := path.Dir(filePath)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(currenDir)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func DbDnsInit() string {
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	database := viper.GetString("mysql.database")

	InitConfig()

	dns := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=utf8&parseTime=True&loc=Local"}, "")

	return dns
}

func RedisURLInit() string {
	password := viper.GetString("redis.password")
	address := viper.GetString("redis.address")
	db := viper.GetString("redis.db")

	InitConfig()
	redisURL := "redis://" + password + "@" + address + "/" + db
	return redisURL
}
