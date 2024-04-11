package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"tiny-tiktok/service/user_service/internal/model"
	"tiny-tiktok/service/user_service/internal/proto"
	"tiny-tiktok/service/user_service/pkg/cache"
	"tiny-tiktok/utils/exceptions"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) UserRegister(ctx context.Context, req *proto.UserRequest) (resp *proto.UserResponse, err error) {
	resp = new(proto.UserResponse)
	var user model.User
	if exist := model.GetInstance().CheckUserNotExist(req.Username); !exist {
		resp.StatusCode = exceptions.UserExist
		resp.StatusMsg = exceptions.GetMsg(exceptions.UserExist)
		return resp, err
	}
	user.UserName = req.Username
	user.PassWord = req.Password

	err = model.GetInstance().Create(&user)
	if err != nil {
		resp.StatusCode = exceptions.DataErr
		resp.StatusMsg = exceptions.GetMsg(exceptions.DataErr)
		return resp, err
	}

	userName, err := model.GetInstance().FindUserByName(req.Username)
	if err != nil {
		resp.StatusCode = exceptions.UserNotExist
		resp.StatusMsg = exceptions.GetMsg(exceptions.UserNotExist)
		return resp, err
	}

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
	resp.UserId = userName.Id
	return resp, nil

}

func (u *UserService) UserLogin(ctx context.Context, req *proto.UserRequest) (resp *proto.UserResponse, err error) {
	resp = new(proto.UserResponse)
	if exist := model.GetInstance().CheckUserNotExist(req.Username); exist {
		resp.StatusCode = exceptions.UserNotExist
		resp.StatusMsg = exceptions.GetMsg(exceptions.UserNotExist)
		return resp, err
	}

	user, err := model.GetInstance().FindUserByName(req.Username)
	if err != nil {
		resp.StatusCode = exceptions.DataErr
		resp.StatusMsg = exceptions.GetMsg(exceptions.DataErr)
		return resp, err
	}
	if ok := model.GetInstance().CheckPassword(req.Password, user.PassWord); !ok {
		resp.StatusCode = exceptions.PasswordError
		resp.StatusMsg = exceptions.GetMsg(exceptions.PasswordError)
		return resp, err
	}

	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)
	resp.UserId = user.Id

	return resp, nil
}

func (u *UserService) UserInfo(ctx context.Context, req *proto.UserInfoRequest) (resp *proto.UserInfoResponse, err error) {
	resp = new(proto.UserInfoResponse)
	userIds := req.UserIds

	for _, userId := range userIds {
		var user *model.User
		key := fmt.Sprintf("%s:%s:%s", "user", "info", strconv.FormatInt(userId, 10))

		exist, err := cache.Redis.Exists(cache.Ctx, key).Result()

		if exist > 0 {
			userString, err := cache.Redis.Get(cache.Ctx, key).Result()
			if err != nil {
				return nil, fmt.Errorf("cache error %v", err)
			}
			err = json.Unmarshal([]byte(userString), &user)
			if err != nil {
				return nil, err
			}
		} else {
			user, err = model.GetInstance().FindUserById(userId)
			if err != nil {
				resp.StatusCode = exceptions.UserNotExist
				resp.StatusMsg = exceptions.GetMsg(exceptions.UserNotExist)
				return resp, err
			}

			userJson, _ := json.Marshal(&user)

			err := cache.Redis.Set(cache.Ctx, key, userJson, 12*time.Hour).Err()
			if err != nil {
				return nil, fmt.Errorf("cache error %v", err)
			}
		}
		resp.Users = append(resp.Users, BuildUser(user))
	}
	resp.StatusCode = exceptions.SUCCESS
	resp.StatusMsg = exceptions.GetMsg(exceptions.SUCCESS)

	return resp, nil
}

func BuildUser(u *model.User) *proto.User {
	user := proto.User{
		Id:              u.Id,
		Name:            u.UserName,
		Avatar:          u.Avatar,
		BackgroundImage: u.BackgroundImage,
		Signature:       u.Signature,
	}
	return &user
}
