package discovery

import (
	"log"
	"net"
	"tiny-tiktok/api_router/pkg/logger"
	"tiny-tiktok/utils/etcd"

	"tiny-tiktok/service/user_service/internal/handler"
	"tiny-tiktok/service/user_service/internal/proto"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func AutoRegister() {
	etcdAddress := viper.GetString("etcd.address")
	etcdRegister, err := etcd.NewServiceRegister([]string{etcdAddress})
	if err != nil {
		logger.Log.Fatal(err)
	}

	serviceName := viper.GetString("server.name")
	serviceAddr := viper.GetString("server.address")
	err = etcdRegister.ServiceRegister(serviceName, serviceAddr, 30)
	if err != nil {
		logger.Log.Fatal(err)
	}

	listener, err := net.Listen("tcp", serviceAddr)
	if err != nil {
		logger.Log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, handler.NewUserService())

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
