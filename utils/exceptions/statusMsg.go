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

	FollowSelfErr: "不能关注自己",
}

func GetMsg(code int) string {
	msg, ok := Msg[code]
	if ok {
		return msg
	}
	return Msg[ERROR]
}
