package model

import (
	"sync"
	"time"
	"tiny-tiktok/utils/snowFlake"
)

type Comment struct {
	Id            int64 `gorm:"primary_key"`
	UserId        int64
	VideoId       int64
	Content       string `gorm:"default:(-)"`
	CreateAt      time.Time
	CommentStatus bool `gorm:"default:(-)"`
}

type CommentModel struct {
}

var commentModel *CommentModel
var commentOnce sync.Once

func GetCommentInstance() *CommentModel {
	commentOnce.Do(
		func() {
			commentModel = &CommentModel{}
		},
	)
	return commentModel
}

func (*CommentModel) AddComment(comment *Comment) (id int64, err error) {
	flake, _ := snowFlake.NewSnowFlake(7, 2)
	comment.Id = flake.NextId()
	comment.CreateAt = time.Now()

	result := DB.Create(&comment)
	if result.Error != nil {
		return -1, result.Error
	}

	err = GetVideoInstance().PlusCommentCount(comment.VideoId)
	if err != nil {
		return -1, err
	}

	return comment.Id, nil
}

func (*CommentModel) DeleteCommentById(commentId int64) error {
	var comment Comment

	result := DB.Table("comment").Where("id=?", commentId).Delete(&comment)

	if result.Error != nil {
		return result.Error
	}

	err := GetVideoInstance().MinusCommentCount(comment.VideoId)
	if err != nil {
		return err
	}

	return nil
}

func (*CommentModel) GetCommentList(videoId int64) ([]Comment, error) {
	var comments []Comment

	result := DB.Table("comment").Where("video_id=?", videoId).Find(&comments)

	if result.Error != nil {
		return nil, result.Error
	}

	return comments, nil
}

func (*CommentModel) GetCommentById(commentId int64) (Comment, error) {
	comment := Comment{}
	result := DB.Table("comment").Where("id=?", commentId).First(&comment)
	if result.Error != nil {
		return comment, result.Error
	}

	return comment, nil
}
