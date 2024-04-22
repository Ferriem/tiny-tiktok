package handler

import (
	"context"
	"fmt"
	"testing"
	"tiny-tiktok/service/user_service/internal/model"
	"tiny-tiktok/service/user_service/internal/proto"
	"tiny-tiktok/service/user_service/pkg/cache"
)

func TestUserInfo(t *testing.T) {
	model.InitDB()
	cache.InitRedis()
	resp, err := NewUserService().UserInfo(context.Background(), &proto.UserInfoRequest{
		UserIds: []int64{249339641491456, 249339641491456, 249339641491456, 249339641491456, 249339641491456},
	})
	fmt.Println(resp, err)
}
