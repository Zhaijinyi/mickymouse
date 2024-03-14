package lock

import (
	"github.com/go-redis/redis"
	"time"
)

func AcquireLock(client *redis.Client, lockKey string, lockValue string, lockDuration time.Duration) (bool, error) {
	lock, err := client.SetNX(lockKey, lockValue, lockDuration).Result()
	if err != nil {
		return false, err
	}
	return lock, nil
}

// 释放分布式锁
func ReleaseLock(client *redis.Client, lockKey string) error {
	_, err := client.Del(lockKey).Result()
	return err
}
