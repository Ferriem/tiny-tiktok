package config

import (
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig() {
	_, filePath, _, _ := runtime.Caller(0)

	currentDir := path.Dir(filePath)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(currentDir)

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

func InitRabbitMQUrl() string {
	user := viper.GetString("rabbitMQ.user")
	password := viper.GetString("rabbitMQ.password")
	address := viper.GetString("rabbitMQ.address")

	url := strings.Join([]string{"amqp://", user, ":", password, "@", address, "/"}, "")
	return url
}
