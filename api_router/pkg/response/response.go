package response

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`

	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`

	TotalFavorited string `json:"total_favorited"`
	WorkCount      int64  `json:"work_count"`
	FavoriteCount  int64  `json:"favorite_count"`
}

type Comment struct {
	Id       int64  `json:"id"`
	User     User   `json:"user"`
	Content  string `json:"content"`
	CreateAt string `json:"create_at"`
}

type Video struct {
	Id            int64  `json:"id"`
	Author        User   `json:"author"`
	VideoUrl      string `json:"video_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

type FeedResponse struct {
	StatusCode int64   `json:"status_code"`
	NextTime   int64   `json:"next_time"`
	StatusMsg  string  `json:"status_msg"`
	VideoList  []Video `json:"video_list"`
}

type UserResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

type UserInfoResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	User       User   `json:"user"`
}

type VideoListResponse struct {
	StatusCode int64   `json:"status_code"`
	StatusMsg  string  `json:"status_msg"`
	VideoList  []Video `json:"video_list"`
}

type PublishActionResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type FavoriteActionResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type CommentActionResponse struct {
	StatusCode int64   `json:"status_code"`
	StatusMsg  string  `json:"status_msg"`
	Comment    Comment `json:"comment"`
}

type CommentDeleteResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type CommentListResponse struct {
	StatusCode  int64     `json:"status_code"`
	StatusMsg   string    `json:"status_msg"`
	CommentList []Comment `json:"comment_list"`
}

type FollowActionResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type FollowListResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserList   []User `json:"user_list"`
}

type PostMessageResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type Message struct {
	Id         int64  `json:"id"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}

type GetMessageResponse struct {
	StatusCode  int64     `json:"status_code"`
	StatusMsg   string    `json:"status_msg"`
	MessageList []Message `json:"message_list"`
}
