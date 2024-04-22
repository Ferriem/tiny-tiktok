package handler

import (
	"context"
	"net/http"
	"strconv"
	"tiny-tiktok/api_router/internal/proto"
	"tiny-tiktok/api_router/pkg/response"
	"tiny-tiktok/utils/exceptions"

	"github.com/gin-gonic/gin"
)

func FavoriteAction(ctx *gin.Context) {
	var favoriteActionReq proto.FavoriteRequest
	userId, _ := ctx.Get("user_id")
	favoriteActionReq.UserId = userId.(int64)
	videoId := ctx.PostForm("video_id")
	if videoId == "" {
		videoId = ctx.Query("video_id")
	}
	favoriteActionReq.VideoId, _ = strconv.ParseInt(videoId, 10, 64)

	actionType := ctx.PostForm("action_type")
	if actionType == "" {
		actionType = ctx.Query("action_type")
	}
	actionTypeValue, _ := strconv.Atoi(actionType)

	if actionTypeValue == 1 || actionTypeValue == 2 {
		favoriteActionReq.ActionType = int64(actionTypeValue)
		videoServiceClient := ctx.Keys["video_service"].(proto.VideoServiceClient)
		videoResp, err := videoServiceClient.Favorite(context.Background(), &favoriteActionReq)
		if err != nil {
			PanicIfFavoriteError(err)
		}

		r := response.FavoriteActionResponse{
			StatusCode: videoResp.StatusCode,
			StatusMsg:  videoResp.StatusMsg,
		}

		ctx.JSON(http.StatusOK, r)
	} else {
		r := response.FavoriteActionResponse{
			StatusCode: exceptions.ErrOperate,
			StatusMsg:  exceptions.GetMsg(exceptions.ErrOperate),
		}
		ctx.JSON(http.StatusOK, r)
	}
}

func FavoriteList(ctx *gin.Context) {
	var favoriteListRequest proto.FavoriteListRequest

	userIdStr := ctx.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	favoriteListRequest.UserId = userId
	videoServiceClient := ctx.Keys["video_service"].(proto.VideoServiceClient)
	favoriteListResp, err := videoServiceClient.FavoriteList(context.Background(), &favoriteListRequest)
	if err != nil {
		PanicIfFavoriteError(err)
	}
	var userIds []int64
	for _, v := range favoriteListResp.VideoList {
		userIds = append(userIds, v.AuthId)
	}

	userInfos := GetUserInfo(userIds, ctx)
	list := BuildVideoList(favoriteListResp.VideoList, userInfos)

	r := response.VideoListResponse{
		StatusCode: favoriteListResp.StatusCode,
		StatusMsg:  favoriteListResp.StatusMsg,
		VideoList:  list,
	}

	ctx.JSON(http.StatusOK, r)

}
