package model

import (
	"sync"
	"time"
	"tiny-tiktok/utils/snowFlake"
)

type Video struct {
	Id       int64 `gorm:"primary_key"`
	AuthId   int64
	Title    string
	CoverUrl string `gorm:"default:(-); unique"`
	VideoUrl string `gorm:"default:(-); unique"`
	CreateAt time.Time
}

type VideoModel struct {
}

var videoModel *VideoModel
var videoOnce sync.Once

func GetVideoInstance() *VideoModel {
	videoOnce.Do(
		func() {
			videoModel = &VideoModel{}
		},
	)
	return videoModel
}

func (*VideoModel) Create(video *Video) error {
	flake, _ := snowFlake.NewSnowFlake(7, 2)
	video.Id = flake.NextId()
	err := DB.Create(&video).Error
	if err != nil {
		return err
	}
	return nil
}

func (*VideoModel) DeleteVideoByUrl(videoUrl string) error {
	var video Video
	if err := DB.Where("video_url=?", videoUrl).First(&video).Error; err != nil {
		return err
	}

	if err := DB.Delete(&video).Error; err != nil {
		return err
	}

	return nil
}

func (*VideoModel) GetVideoByTime(timePoint time.Time) ([]Video, error) {
	var videos []Video
	result := DB.Table("video").Where("create_at < ?", timePoint).Order("create_at DESC").Limit(10).Find(&videos)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(videos) == 0 {
		timePoint = time.Now()
		result = DB.Table("video").Where("create_at < ?", timePoint).Order("create_at DESC").Limit(10).Find(&videos)
		if result.Error != nil {
			return nil, result.Error
		}
		return videos, nil
	}

	return videos, nil
}

func (*VideoModel) GetVideoList(videoIds []int64) ([]Video, error) {
	var videos []Video
	result := DB.Table("video").Where("id IN ?", videoIds).Find(&videos)
	if result.Error != nil {
		return nil, result.Error
	}
	return videos, nil
}

func (*VideoModel) GetVideoListByUser(userId int64) ([]Video, error) {
	var videos []Video
	result := DB.Table("video").Where("auth_id=?", userId).Find(&videos)
	if result.Error != nil {
		return nil, result.Error
	}
	return videos, nil
}

func (*VideoModel) GetFavoriteCount(userId int64) (int64, error) {
	var count int64
	DB.Table("video").Where("auth_id=?", userId).Select("SUM(favorite_count) as count").Pluck("count", &count)
	return count, nil
}

func (*VideoModel) GetWorkCount(userId int64) (int64, error) {
	var count int64
	DB.Table("video").Where("auth_id=?", userId).Count(&count)
	return count, nil
}
