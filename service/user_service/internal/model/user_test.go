package model

import (
	"fmt"
	"testing"
)

func TestUserCreate(t *testing.T) {
	InitDB()
	user := &User{
		UserName: "test",
		PassWord: "123456",
	}
	GetInstance().Create(user)
}

func TestUserFindByName(t *testing.T) {
	InitDB()
	user, _ := GetInstance().FindUserByName("ferriem")
	fmt.Println(user)
}

func TestUserFindById(t *testing.T) {
	InitDB()
	user, _ := GetInstance().FindUserById(565862526663946240)
	fmt.Println(user)
}

func TestUserNotExist(t *testing.T) {
	InitDB()
	exist := GetInstance().CheckUserNotExist("test1")
	fmt.Println(exist)
}
