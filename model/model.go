package model

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Leaderboard struct {
	Client *redis.Client
}

// 创建接收消息的
func SubscribeToChannel(client *redis.Client, channelName string) {
	sub := client.Subscribe(channelName)
	defer sub.Close()   //确保在函数返回前关闭订阅对象，释放资源
	ch := sub.Channel() //获取接受消息的通道
	for msg := range ch {
		fmt.Println("Received message:", msg.Payload)
	}
}

// 创建发布消息的
func PublishToChannel(client *redis.Client, channelName string, message string) error {
	pubsub := client.Publish(channelName, message)
	if pubsub.Err() != nil {
		return pubsub.Err()
	}
	fmt.Println("发布成功!")
	return nil
}
