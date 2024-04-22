package model

import (
	"fmt"
	"sync"
	"time"
	"tiny-tiktok/utils/snowFlake"
)

type Message struct {
	Id       int64 `gorm:"primary_key"`
	UserId   int64
	ToUserId int64
	Content  string
	CreateAt int64
}

type MessageModel struct {
}

var messageModel *MessageModel
var messageOnce sync.Once

func GetMessageInstance() *MessageModel {
	messageOnce.Do(
		func() {
			messageModel = &MessageModel{}
		},
	)
	return messageModel
}

func getCurrentTime() (createTime int64) {
	createTime = time.Now().Unix()
	return
}

func (*MessageModel) PostMessage(message *Message) error {
	message.CreateAt = getCurrentTime()
	flake, _ := snowFlake.NewSnowFlake(7, 3)
	message.Id = flake.NextId()
	err := DB.Create(&message).Error
	return err
}

func (*MessageModel) GetMessage(UserId, ToUserId, PreMsgTime int64, messages *[]Message) error {
	if ToUserId == 0 {
		return fmt.Errorf("ToUserId is 0")
	}
	err := DB.Model(&Message{}).Where(DB.Model(&Message{}).Where(&Message{UserId: UserId, ToUserId: ToUserId}).Or(&Message{UserId: ToUserId, ToUserId: UserId})).Where("create_at > ?", PreMsgTime).Order("create_at").Find(messages).Error
	return err
}
