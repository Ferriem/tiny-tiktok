package discovery

import (
	"log"
	"net"
	"tiny-tiktok/service/video_service/internal/handler"
	"tiny-tiktok/service/video_service/internal/proto"
	"tiny-tiktok/utils/etcd"

	"google.golang.org/grpc"

	"github.com/spf13/viper"
)

func AutoRegister() {
	etcdAddress := viper.GetString("etcd.address")
	etcdRegister, err := etcd.NewServiceRegister([]string{etcdAddress})

	if err != nil {
		log.Fatal(err)
	}

	serviceName := viper.GetString("server.name")
	serviceAddress := viper.GetString("server.address")
	err = etcdRegister.ServiceRegister(serviceName, serviceAddress, 30)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		log.Fatal(err)
	}

	var options = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 128),
	}
	server := grpc.NewServer(options...)
	proto.RegisterVideoServiceServer(server, handler.NewVideoService())

	err = server.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
