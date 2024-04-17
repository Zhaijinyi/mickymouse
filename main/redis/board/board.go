package board

import (
	"github.com/go-redis/redis"
)

// 向排行榜中添加成员及其分数
func AddMemberAndScore(client *redis.Client, member string, score float64) error {
	z := redis.Z{Score: score, Member: member}
	err := client.ZAdd("board", z).Err() //向有序集合中添加一个或多个成员
	return err
}

// 获取排行榜前N名成员
func GetTopMembers(client *redis.Client, n int) ([]string, error) {
	members, err := client.ZRevRange("board", 0, int64(n-1)).Result()//从高到低排序 ZRange从低到高排序
	return members, err
}
