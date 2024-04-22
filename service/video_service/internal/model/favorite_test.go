package model

import (
	"fmt"
	"testing"
)

func TestAddFavorite(t *testing.T) {
	InitDB()
	err := GetFavoriteInstance().Favorite(DB, &Favorite{
		UserId:  1,
		VideoId: 446461502582784,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func TestIsFavorite(t *testing.T) {
	InitDB()
	flag, _ := GetFavoriteInstance().IsFavorite(1, 446461502582784)
	fmt.Println(flag)
}

func TestDeleteFavorite(t *testing.T) {
	InitDB()
	err := GetFavoriteInstance().DeleteFavorite(DB, &Favorite{
		UserId:  1,
		VideoId: 446461502582784,
	})
	if err != nil {
		panic(err)
	}
}

func TestGetFavoriteList(t *testing.T) {
	InitDB()
	var videoId []int64
	videoId, _ = GetFavoriteInstance().FavoriteList(1)
	fmt.Println(videoId)
}

func TestGetFavoriteCount(t *testing.T) {
	InitDB()
	count, _ := GetFavoriteInstance().GetFavoriteCount(1)
	fmt.Println(count)
}

func TestGetVideoFavoriteCount(t *testing.T) {
	InitDB()
	count, _ := GetFavoriteInstance().GetVideoFavoriteCount(1)
	fmt.Println(count)
}

func TestFavoriteUserList(t *testing.T) {
	InitDB()
	userId, _ := GetFavoriteInstance().FavoriteUserList(446461502582784)
	fmt.Println(userId)
}
