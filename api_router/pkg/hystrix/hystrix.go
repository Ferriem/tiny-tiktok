package hystrix

import (
	"errors"
	"tiny-tiktok/api_router/pkg/logger"

	"github.com/afex/hystrix-go/hystrix"
)

func NewWrapper(name string) {
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		Timeout:                1000,
		MaxConcurrentRequests:  10,
		RequestVolumeThreshold: 5,
		SleepWindow:            5000,
		ErrorPercentThreshold:  20,
	})

	err := hystrix.Do(name, func() error {
		return errors.New("熔断")
	}, nil)
	if err != nil {
		logger.Log.Info(name + "熔断")
	}

}
