package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
	"tiny-tiktok/api_router/pkg/logger"
	"tiny-tiktok/service/video_service/internal/model"
	"tiny-tiktok/service/video_service/internal/proto"
	"tiny-tiktok/service/video_service/oss"
	"tiny-tiktok/service/video_service/pkg/cache"
	"tiny-tiktok/service/video_service/pkg/ffmpeg"
	"tiny-tiktok/service/video_service/pkg/mq"
	"tiny-tiktok/utils/exceptions"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type VideoService struct {
	proto.UnimplementedVideoServiceServer
}

func NewVideoService() *VideoService {
	return &VideoService{}
}

func (*VideoService) Feed(ctx context.Context, req *proto.FeedRequest) (resp *proto.FeedResponse, err error) {
	resp = new(proto.FeedResponse)

	var timer time.Time
	if req.LatestTime == -1 {
		timer = time.Now()
	} else {
		timer = time.Unix(req.LatestTime/1000, 0)
	}

	videos, err := model.GetVideoInstance().GetVideoByTime(timer)
	if err != nil {
		resp.StatusCode = exceptions.VideoNoExist
		resp.StatusMsg = exceptions.GetMsg(exceptions.VideoNoExist)
		return resp, err
	}

	if req.UserId == -1 {
		resp.Videos = BuildVideoForFavorite(videos, false)
	} else {
		resp.Videos = BuildVideo(videos, req.UserId)
	}

	LastIndex := len(videos) - 1
	resp.NextTime = videos[LastIndex].CreateAt.Unix()

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)

	return resp, nil
}

func (*VideoService) PublishList(ctx context.Context, req *proto.PublishListRequest) (resp *proto.PublishListResponse, err error) {
	resp = new(proto.PublishListResponse)
	var videos []model.Video

	key := fmt.Sprintf("%s:%s:%s", "user", "work_list", strconv.FormatInt(req.UserId, 10))

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
		videos, err := model.GetVideoInstance().GetVideoListByUser(req.UserId)
		if err != nil {
			resp.StatusCode = exceptions.VideoNoExist
			resp.StatusMsg = exceptions.GetMsg(exceptions.VideoNoExist)
			return resp, err
		}

		videosJson, _ := json.Marshal(videos)
		err = cache.Redis.Set(cache.Ctx, key, videosJson, 30*time.Minute).Err()
		if err != nil {
			return nil, fmt.Errorf("cache error: %v", err)
		}
	}

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
	resp.VideoList = BuildVideo(videos, req.UserId)

	return resp, nil
}

func (*VideoService) CountInfo(ctx context.Context, req *proto.CountRequest) (resp *proto.CountResponse, err error) {
	resp = new(proto.CountResponse)

	userIds := req.UserIds

	for _, userId := range userIds {
		var count proto.Count
		var videos []model.Video

		exists, err := cache.Redis.HExists(cache.Ctx, "user:total_favorite", strconv.FormatInt(userId, 10)).Result()

		if err != nil {
			return nil, fmt.Errorf("cache error: %v", err)
		}

		if !exists {
			totalFavorite := int64(0)
			videos, _ = model.GetVideoInstance().GetVideoListByUser(userId)

			for _, video := range videos {
				videoId := video.Id
				favoriteCount, err := model.GetFavoriteInstance().GetVideoFavoriteCount(videoId)
				if err != nil {
					resp.StatusCode = exceptions.UserNoFavorite
					resp.StatusMsg = exceptions.GetMsg(exceptions.UserNoFavorite)
					return resp, err
				}
				totalFavorite += favoriteCount
			}

			err = cache.Redis.HSet(cache.Ctx, "user:total_favorite", strconv.FormatInt(userId, 10), totalFavorite).Err()
			if err != nil {
				return nil, fmt.Errorf("cache error: %v", err)
			}

			cache.Redis.Expire(cache.Ctx, "user:total_favorite", 30*time.Minute)
		} else {
			count.FavoriteCount, err = cache.Redis.HGet(cache.Ctx, "user:total_favorite", strconv.FormatInt(userId, 10)).Int64()
			if err != nil {
				return nil, fmt.Errorf("cache error: %v", err)
			}
		}

		exists, err = cache.Redis.HExists(cache.Ctx, "user:work_count", strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			return nil, fmt.Errorf("cache error: %v", err)
		}

		if exists {
			count.WorkCount, err = cache.Redis.HGet(cache.Ctx, "user:work_count", strconv.FormatInt(userId, 10)).Int64()
			if err != nil {
				return nil, fmt.Errorf("cache error: %v", err)
			}
		} else {
			count.WorkCount, err = model.GetVideoInstance().GetWorkCount(userId)
			if err != nil {
				resp.StatusCode = exceptions.UserNoVideo
				resp.StatusMsg = exceptions.GetMsg(exceptions.UserNoVideo)
				return resp, err
			}
			err = cache.Redis.HSet(cache.Ctx, "user:work_count", strconv.FormatInt(userId, 10), count.WorkCount).Err()
			if err != nil {
				return nil, fmt.Errorf("cache error: %v", err)
			}
		}

		exists, err = cache.Redis.HExists(cache.Ctx, "user:favorite_count", strconv.FormatInt(userId, 10)).Result()
		if err != nil {
			return nil, fmt.Errorf("cache error: %v", err)
		}

		if exists {
			count.FavoriteCount, err = cache.Redis.HGet(cache.Ctx, "user:favorite_count", strconv.FormatInt(userId, 10)).Int64()
			if err != nil {
				return nil, fmt.Errorf("cache error: %v", err)
			}
		} else {
			count.FavoriteCount, err = model.GetFavoriteInstance().GetFavoriteCount(userId)
			if err != nil {
				resp.StatusCode = exceptions.UserNoFavorite
				resp.StatusMsg = exceptions.GetMsg(exceptions.UserNoFavorite)
				return resp, err
			}
			err = cache.Redis.HSet(cache.Ctx, "user:favorite_count", strconv.FormatInt(userId, 10), count.FavoriteCount).Err()
			if err != nil {
				return nil, fmt.Errorf("cache error: %v", err)
			}
		}

		resp.CountList = append(resp.CountList, &count)
	}

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)

	return resp, nil
}

func PublishVideo() {
	conn := mq.InitMQ()
	ch, err := conn.Channel()
	if err != nil {
		logger.Log.Errorf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"video_publish",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		mq.FailOnError(err, "Failed to declare a queue")
	}

	msgs, err := ch.Consume(
		q.Name,
		"video_service",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		mq.FailOnError(err, "Failed to register a consumer")
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			var req proto.PublishActionRequest
			err := json.Unmarshal(d.Body, &req)
			if err != nil {
				logger.Log.Errorf("Failed to unmarshal: %s", err)
			}
			logger.Log.Infof("upload video from user: %v", req.UserId)
			key := fmt.Sprintf("%s:%s", "user", "work_count")
			title := req.Title
			UUID := uuid.New()
			videoDir := title + "--" + UUID.String() + ".mp4"
			pictureDir := title + "--" + UUID.String() + ".jpg"

			videoUrl := "http://ferriemtiktok.oss-cn-beijing.aliyuncs.com/" + videoDir
			pictureUrl := "http://ferriemtiktok.oss-cn-beijing.aliyuncs.com/" + pictureDir

			var wg sync.WaitGroup
			wg.Add(2)
			var updataErr, createErr error

			go func() {
				defer wg.Done()
				updataErr = oss.UpLoad(videoDir, req.Data)
				logger.Log.Info(updataErr)
				coverByte, _ := ffmpeg.Cover(videoUrl, "00:00:04")
				updataErr = oss.UpLoad(pictureDir, coverByte)
				logger.Log.Info(updataErr)
			}()

			go func() {
				defer wg.Done()
				video := model.Video{
					AuthId:   req.UserId,
					VideoUrl: videoUrl,
					CoverUrl: pictureUrl,
					Title:    req.Title,
					CreateAt: time.Now(),
				}
				createErr = model.GetVideoInstance().Create(&video)
			}()

			if updataErr != nil || createErr != nil {
				go func() {
					if createErr != nil {
						_ = oss.Delete(videoDir)
						_ = oss.Delete(pictureDir)
					}

					if updataErr != nil {
						_ = model.GetVideoInstance().DeleteVideoByUrl(videoUrl)
					}
				}()
			}

			d.Ack(false)

			exist, err := cache.Redis.HExists(cache.Ctx, key, strconv.FormatInt(req.UserId, 10)).Result()
			if err != nil {
				logger.Log.Errorf("cache error: %v", err)
			}
			if exist {
				_, err := cache.Redis.HIncrBy(cache.Ctx, key, strconv.FormatInt(req.UserId, 10), 1).Result()
				if err != nil {
					logger.Log.Errorf("cache error: %v", err)
				}
			}

			workKey := fmt.Sprintf("%s:%s:%s", "user", "work_list", strconv.FormatInt(req.UserId, 10))
			err = cache.Redis.Del(cache.Ctx, workKey).Err()
			if err != nil {
				logger.Log.Errorf("cache error: %v", err)
			}

			go func() {
				time.Sleep(500 * time.Millisecond)
				cache.Redis.Del(cache.Ctx, workKey)
			}()
		}
	}()
	<-forever
}

func BuildVideo(videos []model.Video, userId int64) []*proto.Video {
	var videoResp []*proto.Video

	for _, video := range videos {
		favorite := isFavorite(userId, video.Id)
		favoriteCount := getFavoriteCount(video.Id)
		commentCount := getCommentCount(video.Id)
		videoResp = append(videoResp, &proto.Video{
			Id:            video.Id,
			AuthId:        video.AuthId,
			PlayUrl:       video.VideoUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    favorite,
			Title:         video.Title,
		})
	}

	return videoResp
}

func (*VideoService) PublishAction(ctx context.Context, req *proto.PublishActionRequest) (resp *proto.PublishActionResponse, err error) {
	resp = new(proto.PublishActionResponse)
	reqString, _ := json.Marshal(&req)

	conn := mq.InitMQ()
	ch, err := conn.Channel()
	if err != nil {
		logger.Log.Errorf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"video_publish",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		mq.FailOnError(err, "Failed to declare a queue")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        reqString,
		})
	if err != nil {
		resp.StatusCode = exceptions.VideoUploadErr
		resp.StatusMsg = exceptions.GetMsg(exceptions.VideoUploadErr)
		return resp, err
	}

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)

	return resp, nil
}

func BuildVideoForFavorite(videos []model.Video, isFavorite bool) []*proto.Video {
	var videoResp []*proto.Video

	for _, video := range videos {
		favoriteCount := getFavoriteCount(video.Id)
		commentCount := getCommentCount(video.Id)
		videoResp = append(videoResp, &proto.Video{
			Id:            video.Id,
			AuthId:        video.AuthId,
			PlayUrl:       video.VideoUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         video.Title,
		})
	}

	return videoResp
}
