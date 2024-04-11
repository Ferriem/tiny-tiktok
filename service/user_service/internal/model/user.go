package model

import (
	"sync"
	"tiny-tiktok/service/user_service/pkg/encryption"
	"tiny-tiktok/utils/snowFlake"

	"gorm.io/gorm"
)

type User struct {
	Id              int64  `gorm: "primary_key"`
	UserName        string `gorm: "unique"`
	PassWord        string `gorm: "not null"`
	Avatar          string `gorm:"default:https://ferriem.oss-cn-beijing.aliyuncs.com/106524316_p0_master1200.jpg?Expires=1712740491&OSSAccessKeyId=TMP.3KdzhSjiXUM1fhZajdfXHsFU4iy6AFLL1ipHxo99ghn8aDreboEMSFUjBB6kcbtpyrNZAwsJV1DTC8yjuUodwsN85J4vAz&Signature=2JxF6ODOVLayrAYvQpkRpkIj1B8%3D"`
	BackgroundImage string `gorm:"default:http://ferriem.oss-cn-beijing.aliyuncs.com/91953835_p1.png?Expires=1712741194&OSSAccessKeyId=TMP.3KdzhSjiXUM1fhZajdfXHsFU4iy6AFLL1ipHxo99ghn8aDreboEMSFUjBB6kcbtpyrNZAwsJV1DTC8yjuUodwsN85J4vAz&Signature=iekZXEmhAtdjhSuVigWe08ovQSw%3D"`
	Signature       string `gorm:"default:该用户还没有简介"`
}

type UserModel struct {
}

var userModel *UserModel
var userOnce sync.Once

func GetInstance() *UserModel {
	userOnce.Do(
		func() {
			userModel = &UserModel{}
		},
	)
	return userModel
}

func (*UserModel) Create(user *User) error {
	flake, _ := snowFlake.NewSnowFlake(7, 1)
	user.Id = flake.NextId()
	user.PassWord = encryption.HashPassword(user.PassWord)
	DB.Create(&user)
	return nil
}

func (*UserModel) FindUserByName(username string) (*User, error) {
	user := User{}
	res := DB.Where("user_name=?", username).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (*UserModel) FindUserById(userid int64) (*User, error) {
	user := User{}
	res := DB.Where("id=?", userid).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (*UserModel) CheckUserNotExist(username string) bool {
	user := User{}
	err := DB.Where("user_name=?", username).First(&user).Error
	return err == gorm.ErrRecordNotFound
}

func (*UserModel) CheckPassword(password string, hashPassword string) bool {
	return encryption.VerifyHashPassword(password, hashPassword)
}
