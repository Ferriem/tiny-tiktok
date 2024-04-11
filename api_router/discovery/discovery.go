package discovery

import (
	"tiny-tiktok/api_router/internal/proto"
	"tiny-tiktok/api_router/pkg/logger"
	"tiny-tiktok/utils/etcd"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Resolver() map[string]interface{} {
	serveInstance := make(map[string]interface{})

	etcdAddress := viper.GetString("etcd.address")

	serviceDiscovery, err := etcd.NewServiceDiscovery([]string{etcdAddress})
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer serviceDiscovery.Close()

	// user service
	err = serviceDiscovery.ServiceDiscovery("user_service")
	if err != nil {
		logger.Log.Fatal(err)
	}

	userServiceAddr, _ := serviceDiscovery.GetService("user_service")

	userConn, err := grpc.Dial(
		userServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		logger.Log.Fatal(err)
	}

	userClient := proto.NewUserServiceClient(userConn)

	logger.Log.Info("user service connected")
	serveInstance["user_service"] = userClient

	//social service
	err = serviceDiscovery.ServiceDiscovery("social_service")
	if err != nil {
		logger.Log.Fatal(err)
	}
	socialServiceAddr, _ := serviceDiscovery.GetService("social_service")
	socialConn, err := grpc.Dial(
		socialServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Log.Fatal(err)
	}

	socialClient := proto.NewSocialServiceClient(socialConn)
	logger.Log.Info("social service connected")
	serveInstance["social_service"] = socialClient

	return serveInstance
}
