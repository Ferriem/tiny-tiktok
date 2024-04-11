package snowFlake

import (
	"errors"
	"sync"
	"time"
	"tiny-tiktok/api_router/pkg/logger"
)

type SnowFlake struct {
	lock         sync.Mutex
	timestamp    int64
	workerid     int64
	datacenterid int64
	sequence     int64
}

const (
	epoch             int64 = 1712744936058                                  //2024-04-10 18:28:56
	timestampBits     uint  = 31                                             // 时间戳占用位数
	datacenteridBits  uint  = 5                                              // 数据中心ID占用位数
	worderidBits      uint  = 5                                              // 机器ID占用位数
	sequenceBits      uint  = 12                                             // 序列号占用位数
	timestampMax      int64 = (-1 ^ (-1 << timestampBits))                   // 时间戳最大值
	datacenteridMax   int64 = (-1 ^ (-1 << datacenteridBits))                // 数据中心ID最大值
	worderidMax       int64 = (-1 ^ (-1 << worderidBits))                    // 机器ID最大值
	sequenceMask      int64 = (-1 ^ (-1 << sequenceBits))                    // 序列号掩码
	workeridShift           = sequenceBits                                   // 机器ID左移位数
	datacenteridShift       = sequenceBits + worderidBits                    // 数据中心ID左移位数
	timestampShift          = sequenceBits + worderidBits + datacenteridBits // 时间戳左移位数
)

func NewSnowFlake(workerid, datacenterid int64) (*SnowFlake, error) {
	if workerid > worderidMax || workerid < 0 {
		logger.Log.Errorf("worker Id can't be greater than %d or less than 0", worderidMax)
		return nil, errors.New("worker Id can't be greater than %d or less than 0")
	}
	if datacenterid > datacenteridMax || datacenterid < 0 {
		logger.Log.Errorf("datacenter Id can't be greater than %d or less than 0", datacenteridMax)
		return nil, errors.New("datacenter Id can't be greater than %d or less than 0")
	}
	return &SnowFlake{
		lock:         sync.Mutex{},
		timestamp:    0,
		workerid:     workerid,
		datacenterid: datacenterid,
		sequence:     0,
	}, nil
}

func (s *SnowFlake) NextId() int64 {
	s.lock.Lock()
	now := time.Now().UnixNano() / 1e6
	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		s.lock.Unlock()
		logger.Log.Errorf("epoch must be between 0 and %d", timestampMax-1)
		return 0
	}
	s.timestamp = now
	r := int64((t)<<timestampShift | (s.datacenterid << datacenteridShift) | (s.workerid << workeridShift) | (s.sequence))
	s.lock.Unlock()
	return r
}
