package cache

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var Redis *redis.Client
var ctx = context.Background()

func InitRedis() {
	addr := viper.GetString("redis.address")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")
	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func GenerateFollowKey(UserId int64) string {
	return "Follow:" + strconv.FormatInt(UserId, 10)
}

func GenerateFollowerKey(UserId int64) string {
	return "Follower:" + strconv.FormatInt(UserId, 10)
}

func GenerateFriendKey(UserId int64) string {
	return "Friend:" + strconv.FormatInt(UserId, 10)
}

func IsFollow(UserId int64, ToUserId int64) (bool, error) {
	key := GenerateFollowKey(UserId)
	return Redis.SIsMember(ctx, key, ToUserId).Result()
}

func FollowAction(UserId, ToUserId int64, ActionType int32) error {
	result, err := Redis.SIsMember(ctx, GenerateFollowKey(UserId), ToUserId).Result()
	if err != nil {
		return err
	}

	action := true
	if ActionType == 2 {
		action = false
	}

	if result == action {
		return nil
	}

	pipe := Redis.TxPipeline()
	if action {
		pipe.SAdd(ctx, GenerateFollowKey(UserId), ToUserId)
		pipe.SAdd(ctx, GenerateFollowerKey(ToUserId), UserId)
		_, err := pipe.Exec(ctx)
		if err != nil {
			return err
		}
	} else {
		pipe.SRem(ctx, GenerateFollowKey(UserId), ToUserId)
		pipe.SRem(ctx, GenerateFollowerKey(ToUserId), UserId)
		_, err := pipe.Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetFollowList(reqUser int64, UserId *[]int64) error {
	result, err := Redis.SMembers(ctx, GenerateFollowKey(reqUser)).Result()
	if err != nil {
		return err
	}
	for _, r := range result {
		id, _ := strconv.ParseInt(r, 10, 64)
		*UserId = append(*UserId, id)
	}
	return nil
}

func GetFollowerList(reqUser int64, UserId *[]int64) error {
	result, err := Redis.SMembers(ctx, GenerateFollowerKey(reqUser)).Result()
	if err != nil {
		return err
	}
	for _, r := range result {
		id, _ := strconv.ParseInt(r, 10, 64)
		*UserId = append(*UserId, id)
	}
	return nil
}

func GetFriendList(reqUser int64, UserId *[]int64) error {
	_, err := Redis.Do(ctx, "SINTERSTORE", GenerateFriendKey(reqUser), GenerateFollowKey(reqUser), GenerateFollowerKey(reqUser)).Result()
	if err != nil {
		return err
	}

	result, err := Redis.SMembers(ctx, GenerateFriendKey(reqUser)).Result()
	if err != nil {
		return err
	}

	for _, r := range result {
		id, _ := strconv.ParseInt(r, 10, 64)
		*UserId = append(*UserId, id)
	}
	return nil
}

func GetFollowCount(reqUser int64) (int64, error) {
	return Redis.SCard(ctx, GenerateFollowKey(reqUser)).Result()
}

func GetFollowerCount(reqUser int64) (int64, error) {
	return Redis.SCard(ctx, GenerateFollowerKey(reqUser)).Result()
}
