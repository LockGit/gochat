/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 14:18
 */
package tools

import (
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

var once sync.Once
var RedisClientMap = map[string]*redis.Client{}

type RedisOption struct {
	Host     string
	Port     int
	Password string
	Db       int
}

func GetRedisInstance(redisOpt RedisOption) *redis.Client {
	host := redisOpt.Host
	port := redisOpt.Port
	db := redisOpt.Db
	password := redisOpt.Password
	addr := fmt.Sprintf("%s:%d:%d", host, port, db)
	if redisCli, ok := RedisClientMap[addr]; ok {
		return redisCli
	}
	client := redis.NewClient(&redis.Options{
		Addr:       addr,
		Password:   password,
		DB:         db,
		MaxConnAge: 20 * time.Second,
	})
	RedisClientMap[addr] = client
	return RedisClientMap[addr]
}
