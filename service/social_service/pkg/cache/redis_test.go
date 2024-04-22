package cache

import (
	"fmt"
	"testing"
	"tiny-tiktok/service/social_service/internal/model"
)

func TestFollowAction(t *testing.T) {
	InitRedis()
	err := FollowAction(4, 3, 1)
	fmt.Println(err)
}

func TestIsFollow(t *testing.T) {
	InitRedis()
	isFollow, _ := IsFollow(3, 4)
	fmt.Println(isFollow)
}

func TestGetFollowList(t *testing.T) {
	InitRedis()
	var UserId []int64
	_ = GetFollowList(3, &UserId)
	fmt.Println(UserId)
}

func TestGetFollowerList(t *testing.T) {
	InitRedis()
	var UserId []int64
	_ = GetFollowerList(4, &UserId)
	fmt.Println(UserId)
}

func TestGetFriendList(t *testing.T) {
	InitRedis()
	var UserId []int64
	_ = GetFriendList(3, &UserId)
	fmt.Println(UserId)
}

func TestGetFollowCount(t *testing.T) {
	InitRedis()
	count, _ := GetFollowCount(3)
	fmt.Println(count)
}

func TestGetFollowerCount(t *testing.T) {
	InitRedis()
	count, _ := GetFollowerCount(3)
	fmt.Println(count)
}

func TestAutoSync(t *testing.T) {
	InitRedis()
	model.InitDB()
	AutoSync()
}
