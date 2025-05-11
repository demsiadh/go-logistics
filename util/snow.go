package util

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	workerIDBits     = 5
	datacenterIDBits = 5
	sequenceBits     = 12

	maxWorkerID     = -1 ^ (-1 << workerIDBits)
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits)
	maxSequence     = -1 ^ (-1 << sequenceBits)

	workerIDShift      = sequenceBits
	datacenterIDShift  = sequenceBits + workerIDBits
	timestampLeftShift = sequenceBits + workerIDBits + datacenterIDBits
	sequenceMask       = maxSequence

	// 2020-01-01 00:00:00 UTC 作为时间戳基准点
	twepoch = int64(1577836800000)
)

// Snowflake 雪花算法结构体
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	datacenterID  int64
	sequence      int64
}

// NewSnowflake 创建雪花算法实例
func NewSnowflake(workerID, datacenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("worker ID must be between 0 and 31")
	}
	if datacenterID < 0 || datacenterID > maxDatacenterID {
		return nil, errors.New("datacenter ID must be between 0 and 31")
	}

	return &Snowflake{
		lastTimestamp: -1,
		workerID:      workerID,
		datacenterID:  datacenterID,
		sequence:      0,
	}, nil
}

// NextID 生成下一个唯一ID
func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixNano()/1e6 - twepoch

	if timestamp < s.lastTimestamp {
		return 0, errors.New("clock moved backwards. Refusing to generate id for " + fmt.Sprintf("%d", s.lastTimestamp-timestamp) + " milliseconds")
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 当前毫秒内序列号用完，等待下一毫秒
			timestamp = s.waitNextMillis(s.lastTimestamp)
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	return (timestamp << timestampLeftShift) |
		(s.datacenterID << datacenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence, nil
}

// waitNextMillis 等待下一毫秒
func (s *Snowflake) waitNextMillis(lastTimestamp int64) int64 {
	timestamp := time.Now().UnixNano()/1e6 - twepoch
	for timestamp <= lastTimestamp {
		timestamp = time.Now().UnixNano()/1e6 - twepoch
	}
	return timestamp
}

// GenerateOrderID 生成订单ID
func GenerateOrderID() (string, error) {
	snowflake, err := NewSnowflake(1, 1)
	if err != nil {
		return "", err
	}

	// 生成唯一ID
	id, err := snowflake.NextID()
	if err != nil {
		return "", err
	}

	// 获取当前日期
	currentDate := time.Now().Format("20060102")

	// 提取ID的低6位作为后缀
	suffix := id % 1000000

	// 组合订单号：LY + 日期 + 6位后缀
	orderID := fmt.Sprintf("LY%s%06d", currentDate, suffix)

	return orderID, nil
}
