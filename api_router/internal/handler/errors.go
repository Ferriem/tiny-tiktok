package handler

import (
	"errors"
	"tiny-tiktok/api_router/pkg/logger"
)

func PanicIfUserError(err error) {
	if err != nil {
		err = errors.New("UserService--error" + err.Error())
		logger.Log.Info(err)
		panic(err)
	}
}

func PanicIfSocialError(err error) {
	if err != nil {
		err = errors.New("SocialService--error" + err.Error())
		logger.Log.Info(err)
		panic(err)
	}
}

func PanicIfVideoError(err error) {
	if err != nil {
		err = errors.New("VideoService--error" + err.Error())
		logger.Log.Info(err)
		panic(err)
	}
}

func PanicIfFavoriteError(err error) {
	if err != nil {
		err = errors.New("FavoriteService--error" + err.Error())
		logger.Log.Info(err)
		panic(err)
	}
}

func PanicIfCommentError(err error) {
	if err != nil {
		err = errors.New("CommentService--error" + err.Error())
		logger.Log.Info(err)
		panic(err)
	}
}
