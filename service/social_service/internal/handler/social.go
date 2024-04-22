package handler

import (
	"context"
	"tiny-tiktok/service/social_service/internal/model"
	"tiny-tiktok/service/social_service/internal/proto"
	"tiny-tiktok/service/social_service/pkg/cache"
	"tiny-tiktok/utils/exceptions"
)

type SocialService struct {
	proto.UnimplementedSocialServiceServer
}

func NewSocialService() *SocialService {
	return &SocialService{}
}

func (s *SocialService) FollowAction(ctx context.Context, req *proto.FollowRequest) (resp *proto.FollowResponse, err error) {
	resp = new(proto.FollowResponse)
	if req.UserId == req.ToUserId {
		resp.StatusCode = exceptions.FollowSelfErr
		resp.StatusMsg = exceptions.GetMsg(exceptions.FollowSelfErr)
		return resp, err
	}

	err = cache.FollowAction(req.UserId, req.ToUserId, req.ActionType)
	if err != nil {
		return resp, nil
	}

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
	return resp, nil
}

func (*SocialService) GetFollowList(ctx context.Context, req *proto.FollowListRequest) (resp *proto.FollowListResponse, err error) {
	resp = new(proto.FollowListResponse)
	resp.StatusCode = exceptions.ERROR
	resp.StatusMsg = exceptions.GetMsg(exceptions.ERROR)

	err = cache.GetFollowList(req.UserId, &resp.UserId)

	if err != nil {
		return resp, nil
	}
	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
	return resp, nil
}

func (s *SocialService) GetFollowerList(ctx context.Context, req *proto.FollowListRequest) (resp *proto.FollowListResponse, err error) {
	resp = new(proto.FollowListResponse)
	resp.StatusCode = exceptions.ERROR
	resp.StatusMsg = exceptions.GetMsg(exceptions.ERROR)

	err = cache.GetFollowerList(req.UserId, &resp.UserId)

	if err != nil {
		return resp, nil
	}
	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
	return resp, nil
}

func (s *SocialService) GetFriendList(ctx context.Context, req *proto.FollowListRequest) (resp *proto.FollowListResponse, err error) {
	resp = new(proto.FollowListResponse)
	resp.StatusCode = exceptions.ERROR
	resp.StatusMsg = exceptions.GetMsg(exceptions.ERROR)

	err = cache.GetFriendList(req.UserId, &resp.UserId)

	if err != nil {
		return resp, nil
	}
	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
	return resp, nil
}

func (s *SocialService) GetFollowInfo(ctx context.Context, req *proto.FollowInfoRequest) (resp *proto.FollowInfoResponse, err error) {
	resp = new(proto.FollowInfoResponse)
	for _, toUserId := range req.ToUserId {
		res1, err1 := cache.IsFollow(req.UserId, toUserId)
		cnt2, err2 := cache.GetFollowCount(toUserId)
		cnt3, err3 := cache.GetFollowerCount(toUserId)
		if err1 != nil || err2 != nil || err3 != nil {
			resp.StatusCode = exceptions.ERROR
			resp.StatusMsg = exceptions.GetMsg(exceptions.ERROR)
		}
		resp.FollowInfo = append(resp.FollowInfo, &proto.FollowInfo{
			IsFollow:      res1,
			FollowCount:   cnt2,
			FollowerCount: cnt3,
			ToUserId:      toUserId,
		})
	}
	return resp, nil
}

func (s *SocialService) PostMessage(ctx context.Context, req *proto.PostMessageRequest) (resp *proto.PostMessageResponse, err error) {
	resp = new(proto.PostMessageResponse)
	message := model.Message{
		UserId:   req.UserId,
		ToUserId: req.ToUserId,
		Content:  req.Content,
	}
	err = model.GetMessageInstance().PostMessage(&message)
	if err != nil {
		resp.StatusCode = exceptions.ERROR
		resp.StatusMsg = exceptions.GetMsg(exceptions.ERROR)
		return resp, nil
	}

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
	return resp, nil
}

func (s *SocialService) GetMessage(ctx context.Context, req *proto.GetMessageRequest) (resp *proto.GetMessageResponse, err error) {
	resp = new(proto.GetMessageResponse)
	var messages []model.Message
	err = model.GetMessageInstance().GetMessage(req.UserId, req.ToUserId, req.PreMsgTime, &messages)
	if err != nil {
		resp.StatusCode = exceptions.ERROR
		resp.StatusMsg = exceptions.GetMsg(exceptions.ERROR)
		return resp, nil
	}

	for _, message := range messages {
		resp.Message = append(resp.Message, &proto.Message{
			Id:        message.Id,
			UserId:    message.UserId,
			ToUserId:  message.ToUserId,
			Content:   message.Content,
			CreatedAt: message.CreateAt,
		})
	}

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
	return resp, nil
}
