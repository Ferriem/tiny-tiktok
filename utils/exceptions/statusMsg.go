package exceptions

var Msg = map[int]string{
	SUCCESS: "OK",
	ERROR:   "fail",

	ErrOperate: "操作错误",

	RequestError: "未登录",
	UnAuth:       "Token未授权",
	TokenTimeout: "Token过期",

	DataErr:  "数据错误",
	CacheErr: "Redis异常",

	UserExist:     "用户已存在",
	UserNotExist:  "用户不存在",
	PasswordError: "密码错误",

	VideoNoExist:     "视频不存在",
	VideoUploadErr:   "视频上传失败",
	VideoFavoriteErr: "视频点赞失败",
	UserNoVideo:      "用户无视频",

	FavoriteErr:       "已点赞",
	CancelFavoriteErr: "取消点赞失败",
	VideoNoFavorite:   "视频未点赞",
	UserNoFavorite:    "用户未点赞",

	CommentErr:       "评论失败",
	CommentNoExist:   "评论不存在",
	CommentDeleteErr: "评论删除失败",

	FollowSelfErr: "不能关注自己",
}

func GetMsg(code int) string {
	msg, ok := Msg[code]
	if ok {
		return msg
	}
	return Msg[ERROR]
}
