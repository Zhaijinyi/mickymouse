package board

import (
	"github.com/go-redis/redis"
)

// 向排行榜中添加成员及其分数
func AddMemberAndScore(client *redis.Client, member string, score float64) error {
	z := redis.Z{Score: score, Member: member}
	err := client.ZAdd("board", z).Err() //
	return err
}

// 获取排行榜前N名成员
func GetTopMembers(client *redis.Client, n int) ([]string, error) {
	members, err := client.ZRevRange("board", 0, int64(n-1)).Result()
	return members, err
}
