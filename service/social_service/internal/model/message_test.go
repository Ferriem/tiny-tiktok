package model

import (
	"fmt"
	"testing"
)

func TestPostMessage(t *testing.T) {
	InitDB()
	err := GetMessageInstance().PostMessage(&Message{
		UserId:   1,
		ToUserId: 2,
		Content:  "Hello",
	})
	if err != nil {
		panic(err)
	}
}

func TestGetMessage(t *testing.T) {
	InitDB()
	var messages []Message
	err := GetMessageInstance().GetMessage(2, 1, 0, &messages)
	if err != nil {
		panic(err)
	}
	for _, message := range messages {
		fmt.Println(message)
	}
}
