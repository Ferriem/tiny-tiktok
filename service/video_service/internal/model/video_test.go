package model

import (
	"fmt"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	InitDB()
	Video := Video{
		AuthId:   2,
		Title:    "test",
		VideoUrl: "test1",
		CoverUrl: "test1",
		CreateAt: time.Now(),
	}

	err := GetVideoInstance().Create(&Video)
	fmt.Println(err)
}

func TestDelete(t *testing.T) {
	InitDB()
	err := GetVideoInstance().DeleteVideoByUrl("test")
	fmt.Println(err)
}

func TestGetVideoByTime(t *testing.T) {
	InitDB()
	videos, _ := GetVideoInstance().GetVideoByTime(time.Now())
	fmt.Println(videos)
}

func TestGetVideoList(t *testing.T) {
	InitDB()
	videos, _ := GetVideoInstance().GetVideoList([]int64{446624270938112})
	fmt.Println(videos)
}

func TestGetVideoListByUser(t *testing.T) {
	InitDB()
	videos, _ := GetVideoInstance().GetVideoListByUser(2)
	fmt.Println(videos)
}
