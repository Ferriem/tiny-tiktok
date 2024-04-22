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

func CommentAction(ctx *gin.Context) {
	var commentActionReq proto.CommentRequest

	userId, _ := ctx.Get("user_id")
	commentActionReq.UserId = userId.(int64)

	videoId := ctx.PostForm("video_id")
	if videoId == "" {
		videoId = ctx.Query("video_id")
	}
	commentActionReq.VideoId, _ = strconv.ParseInt(videoId, 10, 64)

	actionType := ctx.PostForm("action_type")
	if actionType == "" {
		actionType = ctx.Query("action_type")
	}

	actionTypeValue, _ := strconv.Atoi(actionType)
	commentActionReq.ActionType = int64(actionTypeValue)

	if commentActionReq.ActionType == 1 {
		content := ctx.PostForm("content")
		if content == "" {
			content = ctx.Query("content")
		}
		commentActionReq.Content = content
	} else if commentActionReq.ActionType == 2 {
		commentId := ctx.PostForm("comment_id")
		if commentId == "" {
			commentId = ctx.Query("comment_id")
		}
		commentActionReq.CommentId, _ = strconv.ParseInt(commentId, 10, 64)
	} else {
		r := response.CommentActionResponse{
			StatusCode: exceptions.ErrOperate,
			StatusMsg:  exceptions.GetMsg(exceptions.ErrOperate),
		}
		ctx.JSON(http.StatusOK, r)
		return
	}

	videoServiceClient := ctx.Keys["video_service"].(proto.VideoServiceClient)
	videoResp, err := videoServiceClient.Comment(context.Background(), &commentActionReq)
	if err != nil {
		PanicIfVideoError(err)
	}

	if actionTypeValue == 1 {
		userIds := []int64{userId.(int64)}
		userInfos := GetUserInfo(userIds, ctx)

		r := response.CommentActionResponse{
			StatusCode: videoResp.StatusCode,
			StatusMsg:  videoResp.StatusMsg,
			Comment:    BuildComment(videoResp.Comment, userInfos[0]),
		}

		ctx.JSON(http.StatusOK, r)
	}
	if actionTypeValue == 2 {
		r := response.CommentActionResponse{
			StatusCode: videoResp.StatusCode,
			StatusMsg:  videoResp.StatusMsg,
		}

		ctx.JSON(http.StatusOK, r)

	}
}

func CommentList(ctx *gin.Context) {
	var commentListReq proto.CommentListRequest

	videoIdStr := ctx.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)
	commentListReq.VideoId = videoId

	videoServiceClient := ctx.Keys["video_service"].(proto.VideoServiceClient)
	commentListResp, err := videoServiceClient.CommentList(context.Background(), &commentListReq)
	if err != nil {
		PanicIfCommentError(err)
	}

	var userIds []int64
	for _, comment := range commentListResp.CommentList {
		userIds = append(userIds, comment.UserId)
	}

	userInfos := GetUserInfo(userIds, ctx)

	commentList := BuildCommentList(commentListResp.CommentList, userInfos)

	r := response.CommentListResponse{
		StatusCode:  commentListResp.StatusCode,
		StatusMsg:   commentListResp.StatusMsg,
		CommentList: commentList,
	}

	ctx.JSON(http.StatusOK, r)
}

func BuildComment(comment *proto.Comment, userInfo response.User) response.Comment {
	return response.Comment{
		Id:       comment.Id,
		Content:  comment.Content,
		CreateAt: comment.CreateAt,
		User:     userInfo,
	}
}

func BuildCommentList(comments []*proto.Comment, userInfos []response.User) []response.Comment {
	var commentList []response.Comment
	for i, comment := range comments {
		commentList = append(commentList, BuildComment(comment, userInfos[i]))
	}
	return commentList
}
