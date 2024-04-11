package model

import (
	"errors"
	"sync"
	"tiny-tiktok/utils/snowFlake"

	"gorm.io/gorm"
)

type Favorite struct {
	Id      int64 `gorm:"primary_key"`
	UserId  int64
	VideoId int64
}

type FavoriteModel struct {
}

var favoriteModel *FavoriteModel
var favoriteOnce sync.Once

func GetFavoriteInstance() *FavoriteModel {
	favoriteOnce.Do(
		func() {
			favoriteModel = &FavoriteModel{}
		},
	)
	return favoriteModel
}

func (*FavoriteModel) AddFavorite(favorite *Favorite) error {
	result := DB.Where("user_id=? AND video_id = ?", favorite.UserId, favorite.VideoId).First(&favorite)

	if result.Error == nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return nil
	} else {
		flake, _ := snowFlake.NewSnowFlake(7, 2)
		favorite.Id = flake.NextId()
		result = DB.Create(&favorite)
		if result.Error != nil {
			return result.Error
		}
		err := GetVideoInstance().PlusFavoriteCount(favorite.VideoId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (*FavoriteModel) IsFavorite(userId, videoId int64) (bool, error) {

	result := DB.Where("user_id=? AND video_id=?", userId, videoId).First(&Favorite{})

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return true, result.Error
	}

	if result.RowsAffected > 0 {
		return true, nil
	}
	return false, nil
}

func (*FavoriteModel) DeleteFavorite(favorite *Favorite) error {
	result := DB.Where("user_id=? AND video_id=?", favorite.UserId, favorite.VideoId).Delete(&favorite)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	if result.RowsAffected > 0 {
		err := GetVideoInstance().MinusFavoriteCount(favorite.VideoId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (*FavoriteModel) GetFavoriteList(userId int64) ([]int64, error) {
	var videoIds []int64
	result := DB.Table("favorite").Where("user_id=?", userId).Pluck("video_id", &videoIds)
	if result.Error != nil {
		return nil, result.Error
	}

	return videoIds, nil
}

func (*FavoriteModel) GetFavoriteCount(userId int64) (int64, error) {
	var count int64
	result := DB.Table("favorite").Where("user_id=?", userId).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (*FavoriteModel) GetVideoFavoriteCount(videoId int64) (int64, error) {
	var count int64
	result := DB.Table("favorite").Where("video_id=?", videoId).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (*FavoriteModel) FavoriteUserList(videoId int64) ([]int64, error) {
	var userIds []int64

	result := DB.Table("favorite").Where("video_id=?", videoId).Pluck("user_id", &userIds)

	if result.Error != nil {
		return nil, result.Error
	}

	return userIds, nil
}
