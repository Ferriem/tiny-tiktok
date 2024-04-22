package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"tiny-tiktok/service/video_service/internal/model"
	"tiny-tiktok/service/video_service/internal/proto"
	"tiny-tiktok/service/video_service/pkg/cache"
	"tiny-tiktok/utils/exceptions"

	"tiny-tiktok/api_router/pkg/logger"

	"github.com/redis/go-redis/v9"
)

func (*VideoService) Comment(ctx context.Context, req *proto.CommentRequest) (resp *proto.CommentResponse, err error) {
	resp = new(proto.CommentResponse)

	key := fmt.Sprintf("%s:%s:%s", "video", "comment_list", strconv.FormatInt(req.VideoId, 10))
	comment := model.Comment{
		VideoId: req.VideoId,
		Content: req.Content,
		UserId:  req.UserId,
	}

	action := req.ActionType

	time := time.Now()

	if action == 1 {
		comment.CreateAt = time

		tx := model.DB.Begin()
		id, err := model.GetCommentInstance().Comment(tx, &comment)
		if err != nil {
			resp.StatusCode = exceptions.CommentErr
			resp.StatusMsg = exceptions.GetMsg(exceptions.CommentErr)
			return resp, err
		}

		comment.Id = id
		comment.CommentStatus = true
		commentJson, _ := json.Marshal(comment)

		member := redis.Z{
			Score:  float64(time.Unix()),
			Member: commentJson,
		}

		err = cache.Redis.ZAdd(cache.Ctx, key, member).Err()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("cache error: %v", err)
		}

		tx.Commit()

		commentResp := &proto.Comment{
			Id:       id,
			Content:  req.Content,
			CreateAt: time.Format("2006-01-02 15:04:05"),
		}

		resp.StatusCode = exceptions.SUCCESS
		resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
		resp.Comment = commentResp

	} else if action == 2 {
		tx := model.DB.Begin()
		commentInstance, _ := model.GetCommentInstance().GetCommentById(tx, req.CommentId)
		commentMarshal, _ := json.Marshal(commentInstance)

		err = model.GetCommentInstance().DeleteCommentById(req.CommentId)
		if err != nil {
			resp.StatusCode = exceptions.CommentDeleteErr
			resp.StatusMsg = exceptions.GetMsg(exceptions.CommentDeleteErr)
			resp.Comment = nil
			return resp, err
		}
		logger.Log.Info(string(commentMarshal))

		count, err := cache.Redis.ZRem(cache.Ctx, key, string(commentMarshal)).Result()
		logger.Log.Info("delete count: ", count)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("cache error: %v", err)
		}

		tx.Commit()

		resp.StatusCode = exceptions.SUCCESS
		resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
		resp.Comment = nil
	}

	return resp, nil
}

func (*VideoService) CommentList(ctx context.Context, req *proto.CommentListRequest) (resp *proto.CommentListResponse, err error) {
	resp = new(proto.CommentListResponse)

	var comments []model.Comment

	key := fmt.Sprintf("%s:%s:%s", "video", "comment_list", strconv.FormatInt(req.VideoId, 10))

	exists, err := cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("cache error: %v", err)
	}

	if exists == 0 {
		err := buildCommentCache(req.VideoId)
		if err != nil {
			return nil, fmt.Errorf("build comment cache error: %v", err)
		}
	}

	commentsString, err := cache.Redis.ZRevRange(cache.Ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("cache error: %v", err)
	}

	for _, commentString := range commentsString {
		var comment model.Comment
		err := json.Unmarshal([]byte(commentString), &comment)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal error: %v", err)
		}
		comments = append(comments, comment)
	}

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
	resp.CommentList = BuildComments(comments)

	return resp, nil
}

func BuildComments(comments []model.Comment) []*proto.Comment {
	var commentList []*proto.Comment
	for _, comment := range comments {
		commentList = append(commentList, &proto.Comment{
			Id:       comment.Id,
			UserId:   comment.UserId,
			Content:  comment.Content,
			CreateAt: comment.CreateAt.Format("2006-01-02 15:04:05"),
		})
	}
	return commentList
}

func buildCommentCache(videoId int64) error {
	key := fmt.Sprintf("%s:%s:%s", "video", "comment_list", strconv.FormatInt(videoId, 10))

	comments, err := model.GetCommentInstance().CommentList(videoId)
	if err != nil {
		return err
	}

	var zMembers []redis.Z
	for _, comment := range comments {
		commentJson, err := json.Marshal(comment)
		if err != nil {
			logger.Log.Errorf("json marshal error: %v", err)
			continue
		}
		zMembers = append(zMembers, redis.Z{
			Score:  float64(comment.CreateAt.Unix()),
			Member: commentJson,
		})
	}

	if len(zMembers) == 0 {
		return nil
	}

	err = cache.Redis.ZAdd(cache.Ctx, key, zMembers...).Err()
	if err != nil {
		return err
	}

	return nil
}

func getCommentCount(videoId int64) int64 {
	key := fmt.Sprintf("%s:%s:%s", "video", "comment_list", strconv.FormatInt(videoId, 10))

	exists, err := cache.Redis.Exists(cache.Ctx, key).Result()
	if err != nil {
		logger.Log.Errorf("cache error: %v", err)
	}

	if exists == 0 {
		err := buildCommentCache(videoId)
		if err != nil {
			logger.Log.Errorf("build comment cache error: %v", err)
		}
	}

	count, err := cache.Redis.ZCard(cache.Ctx, key).Result()
	if err != nil {
		logger.Log.Errorf("cache error: %v", err)
	}

	return count

}
