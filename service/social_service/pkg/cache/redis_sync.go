package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"tiny-tiktok/service/social_service/internal/model"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
)

func getAllKeys(keys *[]string) {
	pattern := "Follow:*"

	iter := Redis.Scan(ctx, 0, pattern, 0).Iterator()

	handleKey := func(key string) error {
		*keys = append(*keys, key)
		return nil
	}

	var g errgroup.Group
	for iter.Next(ctx) {
		key := iter.Val()
		g.Go(func() error {
			return handleKey(key)
		})
	}

	if err := g.Wait(); err != nil {
		panic(err)
	}
}

func getAllValueByKeys(keys []string) string {
	cmds, err := Redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			pipe.SMembers(ctx, key)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	followInfo := make(map[string][]string)

	for index, cmd := range cmds {
		followInfo[keys[index][7:]] = cmd.(*redis.StringSliceCmd).Val()
	}
	marshal, err := json.Marshal(followInfo)
	if err != nil {
		panic(err)
	}

	return string(marshal)
}

func AutoSync() {
	var keys []string
	getAllKeys(&keys)

	res := getAllValueByKeys(keys)
	err := model.GetFollowInstance().RedisToMysql(res)
	if err != nil {
		panic(err)
	}
}

func TimeSync() {
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)), cron.WithLogger(
		cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags)),
	))

	_, err := c.AddFunc("0 */10 * * * *", func() {
		fmt.Println(time.Now(), "start sync")
		AutoSync()
	})

	if err != nil {
		panic(err)
	}

	c.Start()
	select {}
}
