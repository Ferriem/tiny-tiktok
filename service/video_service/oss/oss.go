package oss

import (
	"bytes"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

func getBucket() (*oss.Bucket, error) {
	bucketName := viper.GetString("oss.bucketName")
	endPoint := viper.GetString("oss.endpoint")
	accessKeyID := viper.GetString("oss.accessKeyId")
	accessKeySecret := viper.GetString("oss.accessKeySecret")

	fmt.Println(bucketName, endPoint, accessKeyID, accessKeySecret)
	client, err := oss.New(endPoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func UpLoad(fileDir string, fileBytes []byte) error {
	bucket, err := getBucket()
	if err != nil {
		return err
	}

	file := bytes.NewReader(fileBytes)
	err = bucket.PutObject(fileDir, file)
	return err
}

func Delete(fileDir string) error {
	bucket, err := getBucket()
	if err != nil {
		return err
	}

	err = bucket.DeleteObject(fileDir)
	return err
}
