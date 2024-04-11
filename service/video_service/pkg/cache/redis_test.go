package cache

import (
	"context"
	"testing"
)

func TestInitRedis(t *testing.T) {
	InitRedis()
	Redis.Set(context.Background(), "key2", "value2", 0)
}
