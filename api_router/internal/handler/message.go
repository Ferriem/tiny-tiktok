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

func PostMessage(ctx *gin.Context) {
	var postMessage proto.PostMessageRequest
	userId, _ := ctx.Get("user_id")
	postMessage.UserId = userId.(int64)
	toUserId := ctx.Query("to_user_id")
	postMessage.ToUserId, _ = strconv.ParseInt(toUserId, 10, 64)
	actionType := ctx.Query("action_type")
	actionTypeInt, _ := strconv.ParseInt(actionType, 10, 32)

	if actionTypeInt != 1 {
		r := response.PostMessageResponse{
			StatusCode: exceptions.ErrOperate,
			StatusMsg:  exceptions.GetMsg(exceptions.ErrOperate),
		}
		ctx.JSON(http.StatusOK, r)
		return
	}

	postMessage.ActionType = int32(actionTypeInt)
	content := ctx.Query("content")
	postMessage.Content = content

	socialServiceClient := ctx.Keys["social_service"].(proto.SocialServiceClient)
	socialResp, err := socialServiceClient.PostMessage(context.Background(), &postMessage)
	if err != nil {
		PanicIfSocialError(err)
	}

	r := response.PostMessageResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
	}

	ctx.JSON(http.StatusOK, r)
}

func GetMessage(ctx *gin.Context) {
	var getMessage proto.GetMessageRequest
	userId, _ := ctx.Get("user_id")
	getMessage.UserId = userId.(int64)
	toUserId := ctx.Query("to_user_id")
	getMessage.ToUserId, _ = strconv.ParseInt(toUserId, 10, 64)
	PreMsgTime := ctx.Query("pre_msg_time")
	getMessage.PreMsgTime, _ = strconv.ParseInt(PreMsgTime, 10, 64)

	socialServiceClient := ctx.Keys["social_service"].(proto.SocialServiceClient)
	socialResp, err := socialServiceClient.GetMessage(context.Background(), &getMessage)
	if err != nil {
		PanicIfSocialError(err)
	}

	r := response.GetMessageResponse{
		StatusCode: socialResp.StatusCode,
		StatusMsg:  socialResp.StatusMsg,
	}

	for _, message := range socialResp.Message {
		messageResp := response.Message{
			Id:         message.Id,
			ToUserId:   message.ToUserId,
			FromUserId: message.UserId,
			Content:    message.Content,
			CreateTime: message.CreatedAt,
		}
		r.MessageList = append(r.MessageList, messageResp)
	}

	ctx.JSON(http.StatusOK, r)
}
