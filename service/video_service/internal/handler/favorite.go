package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"tiny-tiktok/api_router/pkg/logger"
	"tiny-tiktok/service/video_service/internal/model"
	"tiny-tiktok/service/video_service/internal/proto"
	"tiny-tiktok/service/video_service/pkg/cache"
	"tiny-tiktok/utils/exceptions"
)

func (*VideoService) Favorite(ctx context.Context, req *proto.FavoriteRequest) (*proto.FavoriteResponse, error) {
	resp := new(proto.FavoriteResponse)
	key := fmt.Sprintf("%s:%s", "user", "favorite_count")
	setKey := fmt.Sprintf("%s:%s:%s", "video", "favorite_video", strconv.FormatInt(req.VideoId, 10))
	favoriteKey := fmt.Sprintf("%s:%s:%s", "user", "favorite_list", strconv.FormatInt(req.UserId, 10))

	action := req.ActionType
	var favorite model.Favorite
	favorite.UserId = req.UserId
	favorite.VideoId = req.VideoId

	setExists, err := cache.Redis.Exists(cache.Ctx, setKey).Result()
	if err != nil {
		return nil, fmt.Errorf("cache error: %v", err)
	}
	if setExists == 0 {
		err := buildVideoFavorite(req.VideoId)
		if err != nil {
			return nil, fmt.Errorf("build video favorite error: %v", err)
		}
	}

	if action == 1 {
		result, err := cache.Redis.SIsMember(cache.Ctx, setKey, req.UserId).Result()
		if err != nil {
			return nil, fmt.Errorf("cache error: %v", err)
		}

		if result {
			resp.StatusCode = exceptions.FavoriteErr
			resp.StatusMsg = exceptions.GetMsg(exceptions.FavoriteErr)
			return resp, err
		}
		tx := model.DB.Begin()
		err = model.GetFavoriteInstance().Favorite(tx, &favorite)
		if err != nil {
			resp.StatusCode = exceptions.FavoriteErr
			resp.StatusMsg = exceptions.GetMsg(exceptions.FavoriteErr)
			return resp, err
		}

		exist, err := cache.Redis.HExists(cache.Ctx, key, strconv.FormatInt(req.UserId, 10)).Result()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("cache error: %v", err)
		}

		if exist {
			_, err = cache.Redis.HIncrBy(cache.Ctx, key, strconv.FormatInt(req.UserId, 10), 1).Result()
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("cache error: %v", err)
			}
		}

		err = cache.Redis.SAdd(cache.Ctx, setKey, strconv.FormatInt(req.UserId, 10)).Err()

		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("cache error: %v", err)
		}

		err = cache.Redis.Del(cache.Ctx, favoriteKey).Err()

		if err != nil {
			return nil, fmt.Errorf("cache error: %v", err)
		}

		defer func() {
			go func() {
				time.Sleep(500 * time.Millisecond)
			}()
		}()

		tx.Commit()
	} else if action == 2 {
		result, err := cache.Redis.SIsMember(cache.Ctx, setKey, req.UserId).Result()
		if err != nil {
			return nil, fmt.Errorf("cache error: %v", err)
		}

		if result == false {
			resp.StatusCode = exceptions.CancelFavoriteErr
			resp.StatusMsg = exceptions.GetMsg(exceptions.FavoriteErr)
			return resp, err
		}

		tx := model.DB.Begin()
		err = model.GetFavoriteInstance().DeleteFavorite(tx, &favorite)
		if err != nil {
			resp.StatusCode = exceptions.CancelFavoriteErr
			resp.StatusMsg = exceptions.GetMsg(exceptions.CancelFavoriteErr)
			return resp, err
		}

		exist, err := cache.Redis.HExists(cache.Ctx, key, strconv.FormatInt(req.UserId, 10)).Result()
		if err != nil {
			return nil, fmt.Errorf("cache error: %v", err)
		}

		if exist {
			_, err = cache.Redis.HIncrBy(cache.Ctx, key, strconv.FormatInt(req.UserId, 10), -1).Result()
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("cache error: %v", err)
			}
		}

		err = cache.Redis.SRem(cache.Ctx, setKey, strconv.FormatInt(req.UserId, 10)).Err()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("cache error: %v", err)
		}

		err = cache.Redis.Del(cache.Ctx, favoriteKey).Err()
		if err != nil {
			return nil, fmt.Errorf("cache error: %v", err)
		}

		defer func() {
			go func() {
				time.Sleep(500 * time.Millisecond)
			}()
		}()

		tx.Commit()
	}

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)

	return resp, nil
}

func (*VideoService) FavoriteList(ctx context.Context, req *proto.FavoriteListRequest) (resp *proto.FavoriteListResponse, err error) {
	resp = new(proto.FavoriteListResponse)
	var videos []model.Video
	key := fmt.Sprintf("%s:%s:%s", "user", "favorite_list", strconv.FormatInt(req.UserId, 10))

	exists, err := cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("cache error: %v", err)
	}

	if exists > 0 {
		videosString, err := cache.Redis.Get(cache.Ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("cache error: %v", err)
		}
		err = json.Unmarshal([]byte(videosString), &videos)
		if err != nil {
			return nil, err
		}
	} else {
		var videoIds []int64
		videoIds, err = model.GetFavoriteInstance().FavoriteList(req.UserId)
		if err != nil {
			resp.StatusCode = exceptions.UserNoVideo
			resp.StatusMsg = exceptions.GetMsg(exceptions.UserNoVideo)
			return resp, err
		}

		videos, err = model.GetVideoInstance().GetVideoList(videoIds)
		if err != nil {
			resp.StatusCode = exceptions.VideoNoExist
			resp.StatusMsg = exceptions.GetMsg(exceptions.VideoNoExist)
			return resp, err
		}

		videosJson, _ := json.Marshal(videos)
		err := cache.Redis.Set(cache.Ctx, key, videosJson, 30*time.Minute).Err()
		if err != nil {
			return nil, fmt.Errorf("cache error: %v", err)
		}
	}

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
	resp.VideoList = BuildVideoForFavorite(videos, true)

	return resp, nil
}

func isFavorite(userId int64, videoId int64) bool {
	var isFavorite bool
	key := fmt.Sprintf("%s:%s:%s", "video", "favorite_video", strconv.FormatInt(videoId, 10))

	exists, err := cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		logger.Log.Errorf("cache error: %v", err)
	}

	if exists > 0 {
		isFavorite, err = cache.Redis.SIsMember(cache.Ctx, key, strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			logger.Log.Errorf("cache error: %v", err)
		}
	} else {
		err := buildVideoFavorite(videoId)
		if err != nil {
			logger.Log.Errorf("build video favorite error: %v", err)
		}
		isFavorite, err = cache.Redis.SIsMember(cache.Ctx, key, strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			logger.Log.Errorf("cache error: %v", err)
		}
	}

	return isFavorite
}

func buildVideoFavorite(videoId int64) error {
	key := fmt.Sprintf("%s:%s:%s", "video", "favorite_video", strconv.FormatInt(videoId, 10))

	userIdList, err := model.GetFavoriteInstance().FavoriteUserList(videoId)

	if err != nil {
		return err
	}

	userIds := make([]interface{}, len(userIdList))
	for i, v := range userIdList {
		userIds[i] = v
	}

	if len(userIds) == 0 {
		return nil
	}
	err = cache.Redis.SAdd(cache.Ctx, key, userIds...).Err()
	if err != nil {
		return err
	}

	return nil
}

func getFavoriteCount(videoId int64) int64 {
	setKey := fmt.Sprintf("%s:%s:%s", "video", "favorite_video", strconv.FormatInt(videoId, 10))

	setExists, err := cache.Redis.Exists(cache.Ctx, setKey).Result()
	if err != nil {
		logger.Log.Errorf("cache error: %v", err)
	}

	if setExists == 0 {
		err := buildVideoFavorite(videoId)
		if err != nil {
			logger.Log.Errorf("build video favorite error: %v", err)
		}
	}

	count, err := cache.Redis.SCard(cache.Ctx, setKey).Result()

	if err != nil {
		logger.Log.Errorf("cache error: %v", err)
	}

	return count
}
