package model

import (
	"fmt"
	"testing"
)

func TestFollowAction(t *testing.T) {
	InitDB()
	follow := &Follow{
		UserId:   2,
		ToUserId: 1,
		IsFollow: 1,
	}
	err := GetFollowInstance().FollowAction(follow)
	fmt.Println(err)
}

func TestIsFollow(t *testing.T) {
	InitDB()
	isFollow, _ := GetFollowInstance().IsFollow(1, 2)
	fmt.Println(isFollow)
}

func TestGetFollowList(t *testing.T) {
	InitDB()
	var UserId []int64
	_ = GetFollowInstance().GetFollowList(1, &UserId)
	fmt.Println(UserId)
}

func TestGetFollowerList(t *testing.T) {
	InitDB()
	var UserId []int64
	_ = GetFollowInstance().GetFollowerList(2, &UserId)
	fmt.Println(UserId)
}

func TestGetFriendList(t *testing.T) {
	InitDB()
	var UserId []int64
	_ = GetFollowInstance().GetFriendList(1, &UserId)
	fmt.Println(UserId)

}

func TestGetFollowCount(t *testing.T) {
	InitDB()
	count, _ := GetFollowInstance().GetFollowCount(1)
	fmt.Println(count)
}

func TestGetFollowerCount(t *testing.T) {
	InitDB()
	count, _ := GetFollowInstance().GetFollowerCount(2)
	fmt.Println(count)
}
