package handler

import (
	"context"
	"io"
	"strconv"
	"tiny-tiktok/api_router/internal/proto"
	"tiny-tiktok/api_router/pkg/logger"
	"tiny-tiktok/api_router/pkg/response"

	"github.com/gin-gonic/gin"
)

func Feed(ctx *gin.Context) {
	var feedReq proto.FeedRequest

	token := ctx.Query("token")
	if token == "" {
		feedReq.UserId = -1
	} else {
		userId := ctx.Query("user_id")
		feedReq.UserId, _ = strconv.ParseInt(userId, 10, 64)
	}

	latestTime := ctx.Query("latest_time")
	if latestTime == "" || latestTime == "0" {
		feedReq.LatestTime = -1
	} else {
		timePoint, _ := strconv.ParseInt(latestTime, 10, 64)
		feedReq.LatestTime = timePoint
	}
	videoServiceClient := ctx.Keys["video_service"].(proto.VideoServiceClient)
	feedResp, err := videoServiceClient.Feed(context.Background(), &feedReq)
	if err != nil {
		PanicIfVideoError(err)
	}

	var userIds []int64
	for _, video := range feedResp.Videos {
		userIds = append(userIds, video.AuthId)
	}
	logger.Log.Info(userIds)

	userInfos := GetUserInfo(userIds, ctx)

	list := BuildVideoList(feedResp.Videos, userInfos)

	r := response.FeedResponse{
		StatusCode: feedResp.StatusCode,
		StatusMsg:  feedResp.StatusMsg,
		VideoList:  list,
		NextTime:   feedResp.NextTime,
	}

	ctx.JSON(200, r)
}

func PublishAction(ctx *gin.Context) {
	var publishReq proto.PublishActionRequest

	userId, _ := ctx.Get("user_id")
	publishReq.UserId = userId.(int64)

	title := ctx.PostForm("title")
	publishReq.Title = title

	formFile, _ := ctx.FormFile("data")
	file, err := formFile.Open()
	if err != nil {
		PanicIfVideoError(err)
	}
	defer file.Close()
	buf, err := io.ReadAll(file)
	if err != nil {
		PanicIfVideoError(err)
	}
	publishReq.Data = buf

	videoServiceClient := ctx.Keys["video_service"].(proto.VideoServiceClient)
	videoServiceResp, err := videoServiceClient.PublishAction(context.Background(), &publishReq)
	if err != nil {
		PanicIfVideoError(err)
	}

	r := response.PublishActionResponse{
		StatusCode: videoServiceResp.StatusCode,
		StatusMsg:  videoServiceResp.StatusMsg,
	}

	ctx.JSON(200, r)
}

func PublishList(ctx *gin.Context) {
	var publishListRequest proto.PublishListRequest

	userIdStr := ctx.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)

	publishListRequest.UserId = userId

	videoServiceClient := ctx.Keys["video_service"].(proto.VideoServiceClient)
	publishListResp, err := videoServiceClient.PublishList(context.Background(), &publishListRequest)

	if err != nil {
		PanicIfVideoError(err)
	}

	var userIds []int64
	for _, v := range publishListResp.VideoList {
		userIds = append(userIds, v.AuthId)
	}

	logger.Log.Info(userIds)

	userInfos := GetUserInfo(userIds, ctx)

	list := BuildVideoList(publishListResp.VideoList, userInfos)

	r := response.VideoListResponse{
		StatusCode: publishListResp.StatusCode,
		StatusMsg:  publishListResp.StatusMsg,
		VideoList:  list,
	}

	ctx.JSON(200, r)
}

func BuildVideoList(videos []*proto.Video, userInfos []response.User) []response.Video {
	var list []response.Video
	for i, video := range videos {
		list = append(list, response.Video{
			Id:            video.Id,
			Author:        userInfos[i],
			VideoUrl:      video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		})
	}
	return list
}
