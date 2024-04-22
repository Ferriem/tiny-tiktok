package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"tiny-tiktok/api_router/internal/proto"
	"tiny-tiktok/api_router/pkg/auth"
	"tiny-tiktok/api_router/pkg/response"
	"tiny-tiktok/utils/exceptions"

	"github.com/gin-gonic/gin"
)

func UserRegister(ctx *gin.Context) {
	var userReq proto.UserRequest

	userName := ctx.PostForm("username")
	if userName == "" {
		userName = ctx.Query("username")
	}
	userReq.Username = userName

	passWord := ctx.PostForm("password")
	if passWord == "" {
		passWord = ctx.Query("password")
	}
	userReq.Password = passWord

	userServiceClient := ctx.Keys["user_service"].(proto.UserServiceClient)
	userResp, err := userServiceClient.UserRegister(context.Background(), &userReq)
	if err != nil {
		PanicIfUserError(err)
	}

	token := ""
	if userResp.UserId != 0 {
		token, _ = auth.GenerateToken(userResp.UserId)
	}

	r := response.UserResponse{
		StatusCode: userResp.StatusCode,
		StatusMsg:  userResp.StatusMsg,
		UserId:     userResp.UserId,
		Token:      token,
	}

	ctx.JSON(http.StatusOK, r)
}

func UserLogin(ctx *gin.Context) {
	var userReq proto.UserRequest
	userName := ctx.PostForm("username")
	if userName == "" {
		userName = ctx.Query("username")
	}
	userReq.Username = userName

	passWord := ctx.PostForm("password")
	if passWord == "" {
		passWord = ctx.Query("password")
	}
	userReq.Password = passWord

	userServiceClient := ctx.Keys["user_service"].(proto.UserServiceClient)
	userResp, err := userServiceClient.UserLogin(context.Background(), &userReq)
	if err != nil {
		PanicIfUserError(err)
	}

	token := ""
	if userResp.UserId != 0 {
		token, _ = auth.GenerateToken(userResp.UserId)
	}

	r := response.UserResponse{
		StatusCode: userResp.StatusCode,
		StatusMsg:  userResp.StatusMsg,
		UserId:     userResp.UserId,
		Token:      token,
	}

	ctx.JSON(http.StatusOK, r)
}

func UserInfo(ctx *gin.Context) {
	var userIds []int64

	userIdStr := ctx.Query("user_id")
	userId, _ := strconv.ParseInt(userIdStr, 10, 64)
	userIds = append(userIds, userId)

	r := response.UserInfoResponse{
		StatusCode: exceptions.SUCCESS,
		StatusMsg:  exceptions.GetMsg(exceptions.SUCCESS),
		User:       GetUserInfo(userIds, ctx)[0],
	}
	ctx.JSON(http.StatusOK, r)
}

func GetUserInfo(userIds []int64, ctx *gin.Context) (userInfos []response.User) {
	var err error

	var userInfoReq proto.UserInfoRequest
	var countInfoReq proto.CountRequest
	var followInfoReq proto.FollowInfoRequest

	userInfoReq.UserIds = userIds
	countInfoReq.UserIds = userIds
	followInfoReq.ToUserId = userIds

	var userResp *proto.UserInfoResponse
	var countInfoResp *proto.CountResponse
	var followInfoResp *proto.FollowInfoResponse

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		userServiceClient := ctx.Keys["user_service"].(proto.UserServiceClient)
		userResp, err = userServiceClient.UserInfo(context.Background(), &userInfoReq)
		if err != nil {
			PanicIfUserError(err)
		}
	}()

	go func() {
		defer wg.Done()
		videoServiceClient := ctx.Keys["video_service"].(proto.VideoServiceClient)
		countInfoResp, err = videoServiceClient.CountInfo(context.Background(), &countInfoReq)
		if err != nil {
			PanicIfVideoError(err)
		}
	}()

	go func() {
		defer wg.Done()
		userIdStr, _ := ctx.Get("user_id")
		userId, _ := userIdStr.(int64)
		followInfoReq.UserId = userId
		socialServiceClient := ctx.Keys["social_service"].(proto.SocialServiceClient)
		followInfoResp, err = socialServiceClient.GetFollowInfo(context.Background(), &followInfoReq)
		if err != nil {
			PanicIfSocialError(err)
		}
	}()

	wg.Wait()
	fmt.Println(userIds)
	for index := range userIds {
		userInfos = append(userInfos, BuildUser(userResp.Users[index], countInfoResp.CountList[index], followInfoResp.FollowInfo[index]))
	}
	return userInfos
}

func BuildUser(user *proto.User, count *proto.Count, follow *proto.FollowInfo) response.User {
	return response.User{
		Id:   user.Id,
		Name: user.Name,

		FollowCount:   follow.FollowCount,
		FollowerCount: follow.FollowerCount,
		IsFollow:      follow.IsFollow,

		TotalFavorited: strconv.FormatInt(count.FavoriteCount, 10),
		WorkCount:      count.WorkCount,
		FavoriteCount:  count.FavoriteCount,

		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
	}
}
