package lock

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"log"
	"sync"
	"time"
)

// 全局客户实例 避免多次重复创建
var RedisClient *redis.Client

// initRedisClient 初始化RedisClient
func InitRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

// AcquireLock 尝试获取分布式锁，并设置自动续期功能
func AcquireLock(lockKey, lockValue string, lockDuration time.Duration) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background()) //创建了一个可以被取消的上下文，并传递给Redis命令，以确保在锁不再需要时能够正确清理资源
	var wg sync.WaitGroup                                   //确保在返回cancel函数之前，自动续期的goroutine已经启动并准备执行
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(lockDuration / 2) // 自动间隔设为锁时长的一半
		defer ticker.Stop()
		lock, err := RedisClient.SetNX(lockKey, lockValue, lockDuration).Result() //SETNX命令尝试设置一个键值对,解决锁的粒度问题，确保不存在重复的锁 锁的粒度通过使用唯一的锁标识（lockKey和lockValue）来管理，确保不同的操作或资源被不同的锁控制
		if err != nil {
			log.Printf("Failed to acquire lock for key %s: %v", lockKey, err)
			return
		}
		if !lock {
			log.Printf("Lock for key %s is already held by another client", lockKey)
			return
		}
		for {
			select {
			case <-ctx.Done():
				// 如果外部调用了cancel，则释放锁并退出goroutine
				if err := ReleaseLock(lockKey); err != nil {
					log.Printf("Failed to release lock for key %s: %v", lockKey, err)
				}
				return
			case <-ticker.C:
				// 定期续期锁
				if _, err := RedisClient.Expire(lockKey, lockDuration).Result(); err != nil {
					log.Printf("Failed to renew lock for key %s: %v", lockKey, err)
					return
				}
				log.Printf("Lock for key %s renewed", lockKey)
			}
		}
	}()
	// 等待goroutine启动完成（确保续期逻辑开始执行）
	wg.Wait()
	return cancel, nil
}

// ReleaseLock 释放分布式锁
func ReleaseLock(lockKey string) error {
	ctx := "context..."
	// 调用 Del 方法，传递上下文和键列表（这里只有一个键）
	_, err := RedisClient.Del(ctx, lockKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) { // 键不存在视为释放成功
			log.Printf("Lock for key %s does not exist, assuming it's already released", lockKey)
			return nil
		}
		return err
	}
	log.Printf("Lock for key %s released", lockKey)
	return nil
}
