package exceptions

//* 错误码：
//* 四位组成
//* 1. 1开头代表用户端错误
//* 2. 2开头代表当前系统异常
//* 3. 3开头代表第三方服务异常
//* 4. 4开头若无法确定具体错误，选择宏观错误

const (
	SUCCESS = 0
	ERROR   = -1

	RequestError = 1000 //token
	UnAuth       = 1001
	TokenTimeout = 1002

	ErrOperate = 1100

	DataErr  = 2000
	CacheErr = 2001

	UserExist     = 2100
	UserNotExist  = 2101
	PasswordError = 2102

	VideoNoExist     = 2200
	VideoUploadErr   = 2201
	VideoFavoriteErr = 2203
	UserNoVideo      = 2204

	FavoriteErr       = 2300
	CancelFavoriteErr = 2301
	VideoNoFavorite   = 2302
	UserNoFavorite    = 2303

	CommentErr       = 2400
	CommentNoExist   = 2401
	CommentDeleteErr = 2402

	FollowSelfErr = 2500
)
