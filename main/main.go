package main

import (
	"awesomeProject14/board"
	"awesomeProject14/lock"
	"awesomeProject14/model"
	"log"
	"time"
)

func main() {
	// 初始化Redis客户端
	lock.InitRedisClient()
	client := lock.RedisClient // 全局RedisClient实例
	member := "student1"
	score := 100.0
	n := 10
	channelname := "micky"
	message := "4U"
	lockkey := "mylock"
	lockvalue := "unique" // 唯一标识符
	lockduration := 30 * time.Second
	board.AddMemberAndScore(client, member, score)
	board.GetTopMembers(client, n)
	model.SubscribeToChannel(client, channelname)
	model.PublishToChannel(client, channelname, message)
	cancel, err := lock.AcquireLock(lockkey, lockvalue, lockduration)
	if err != nil {
		log.Fatalf("Failed to acquire lock: %v", err) //日志记录错误
	}
	defer cancel() // 使用defer来确保取消续期并尝试释放锁，相较于之前显式调用释放函数有所修改
}
