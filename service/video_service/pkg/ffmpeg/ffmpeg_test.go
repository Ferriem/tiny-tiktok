package ffmpeg

import (
	"testing"
	"tiny-tiktok/service/video_service/config"
	"tiny-tiktok/service/video_service/oss"
)

func TestCover(t *testing.T) {
	config.InitConfig()
	videoURL := "http://ferriemtiktok.oss-cn-beijing.aliyuncs.com/若能绽放光芒.flv"
	imageBytes, _ := Cover(videoURL, "00:00:04")
	oss.UpLoad("cover.jpg", imageBytes)
}
