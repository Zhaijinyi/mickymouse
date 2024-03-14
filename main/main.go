package main

import (
	"awesomeProject11/main/redis/board"
	"awesomeProject11/main/redis/board/model"
	"awesomeProject11/main/redis/board/model/lock"
	"github.com/go-redis/redis"
	"time"
)

func main() {
	// 连接到Redis数据库
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	member := "student1"
	score := 100.0
	n := 10
	channelname := "micky"
	message := "4U"
	lockkey := "mylock"
	lockvalue := "locked"
	lockduration := 30 * time.Second
	board.AddMemberAndScore(client, member, score)
	board.GetTopMembers(client, n)
	model.SubscribeToChannel(client, channelname)
	model.PublishToChannel(client, channelname, message)
	lock.AcquireLock(client, lockkey, lockvalue, lockduration)
	lock.ReleaseLock(client, lockkey)
}
