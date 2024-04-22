package handler

import (
	"context"
	"fmt"
	"testing"
	"tiny-tiktok/service/video_service/internal/model"
	"tiny-tiktok/service/video_service/internal/proto"
	"tiny-tiktok/service/video_service/pkg/cache"
)

func TestBuildVideoFavorite(t *testing.T) {
	model.InitDB()
	cache.InitRedis()
	err := buildVideoFavorite(1820190535086080)
	fmt.Println(err)
}

func TestGetFavoriteCount(t *testing.T) {
	model.InitDB()
	cache.InitRedis()
	count := getFavoriteCount(1820190535086080)
	fmt.Println(count)
}

func TestFavoriteAction(t *testing.T) {
	model.InitDB()
	cache.InitRedis()
	resp, err := NewVideoService().Favorite(context.Background(), &proto.FavoriteRequest{
		UserId:     1138556563456,
		VideoId:    1820190535086080,
		ActionType: 1,
	})
	fmt.Println(resp, err)
}
