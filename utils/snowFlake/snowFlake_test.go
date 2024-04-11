package snowFlake

import (
	"fmt"
	"testing"
)

func TestSnowFlake_NextId(t *testing.T) {
	snow, _ := NewSnowFlake(10, 10)
	for i := 0; i < 10; i++ {
		fmt.Println(snow.NextId())
	}
}
