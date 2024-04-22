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

type Follow struct {
	IsFollow      bool
	FollowCount   int64
	FollowerCount int64
}

func FollowAction(ctx *gin.Context) {
	var followAction proto.FollowRequest
	userId, _ := ctx.Get("user_id")
	followAction.UserId = userId.(int64)
	toUserId := ctx.Query("to_user_id")
	followAction.ToUserId, _ = strconv.ParseInt(toUserId, 10, 64)
	actionType := ctx.Query("action_type")
	actionTypeInt, _ := strconv.ParseInt(actionType, 10, 32)
	followAction.ActionType = int32(actionTypeInt)

	if actionTypeInt != 1 && actionTypeInt != 2 {
		r := response.FollowActionResponse{
			StatusCode: exceptions.ErrOperate,
			StatusMsg:  exceptions.GetMsg(exceptions.ErrOperate),
		}
		ctx.JSON(http.StatusOK, r)
		return
	}

	socialServiceClient := ctx.Keys["social_service"].(proto.SocialServiceClient)
	socialResp, err := socialServiceClient.FollowAction(context.Background(), &followAction)
	if err != nil {
		PanicIfSocialError(err)
	}

	r := response.FollowActionResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
	}

	ctx.JSON(http.StatusOK, r)
}

func GetFollowList(ctx *gin.Context) {
	var followListRequest proto.FollowListRequest
	userId, _ := ctx.Get("user_id")
	followListRequest.UserId = userId.(int64)

	socialServiceClient := ctx.Keys["social_service"].(proto.SocialServiceClient)
	socialResp, err := socialServiceClient.GetFollowList(context.Background(), &followListRequest)
	if err != nil {
		PanicIfSocialError(err)
	}

	r := response.FollowListResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
		UserList:   GetUserInfo(socialResp.UserId, ctx),
	}

	ctx.JSON(http.StatusOK, r)

}

func GetFollowerList(ctx *gin.Context) {
	var followerListRequest proto.FollowListRequest
	userId, _ := ctx.Get("user_id")
	followerListRequest.UserId = userId.(int64)

	socialServiceClient := ctx.Keys["social_service"].(proto.SocialServiceClient)
	socialResp, err := socialServiceClient.GetFollowerList(context.Background(), &followerListRequest)
	if err != nil {
		PanicIfSocialError(err)
	}

	r := response.FollowListResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
		UserList:   GetUserInfo(socialResp.UserId, ctx),
	}

	ctx.JSON(http.StatusOK, r)
}

func GetFriendList(ctx *gin.Context) {
	var friendListRequest proto.FollowListRequest
	userId, _ := ctx.Get("user_id")
	friendListRequest.UserId = userId.(int64)

	socialServiceClient := ctx.Keys["social_service"].(proto.SocialServiceClient)
	socialResp, err := socialServiceClient.GetFriendList(context.Background(), &friendListRequest)
	if err != nil {
		PanicIfSocialError(err)
	}

	r := response.FollowListResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
		UserList:   GetUserInfo(socialResp.UserId, ctx),
	}

	ctx.JSON(http.StatusOK, r)
}
