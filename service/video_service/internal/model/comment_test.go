package model

import (
	"fmt"
	"testing"
)

func TestAddComment(t *testing.T) {
	InitDB()
	comment := &Comment{
		VideoId: 446461502582784,
		Content: "test",
	}
	id, _ := GetCommentInstance().AddComment(comment)
	fmt.Println(id)
}

func TestDeleteComment(t *testing.T) {
	InitDB()
	err := GetCommentInstance().DeleteCommentById(463969626386432)
	fmt.Println(err)
}

func TestGetCommentList(t *testing.T) {
	InitDB()
	var comments []Comment
	comments, _ = GetCommentInstance().GetCommentList(446461502582784)
	fmt.Println(comments)
}

func TestGetCommentById(t *testing.T) {
	InitDB()
	comment, _ := GetCommentInstance().GetCommentById(463829893148672)
	fmt.Println(comment)
}
