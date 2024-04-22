package handler

import (
	"context"
	"fmt"
	"testing"
	"tiny-tiktok/service/video_service/internal/model"
	"tiny-tiktok/service/video_service/internal/proto"
	"tiny-tiktok/service/video_service/pkg/cache"
)

func TestFeed(t *testing.T) {
	model.InitDB()
	cache.InitRedis()
	resp, err := NewVideoService().Feed(context.Background(), &proto.FeedRequest{
		UserId:     1138556563456,
		LatestTime: -1,
	})
	fmt.Println(resp, err)
}

func TestCountInfo(t *testing.T) {
	model.InitDB()
	cache.InitRedis()
	resp, err := NewVideoService().CountInfo(context.Background(), &proto.CountRequest{
		UserIds: []int64{249339641491456, 249339641491456, 249339641491456, 249339641491456, 249339641491456, 2, 2},
	})
	fmt.Println(resp, err)
}

func Test(t *testing.T) {
	model.InitDB()
	cache.InitRedis()

	count, err := cache.Redis.ZRem(cache.Ctx, "video:comment_list:1816632293093376", "{\"Id\":2220318735495168,\"UserId\":1138556563456,\"VideoId\":1816632293093376,\"Content\":\"\xe5\xa5\xbd\",\"CreateAt\":\"2024-04-16T21:31:41.296962+08:00\",\"CommentStatus\":true}").Result()
	fmt.Println(count, err)
}
