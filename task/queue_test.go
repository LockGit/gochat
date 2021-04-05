/**
 * Created by lock
 * Date: 2021/4/5
 */
package task

import (
	"gochat/config"
	"gochat/tools"
	"testing"
	"time"
)

func Test_TestQueue(t *testing.T) {
	redisOpt := tools.RedisOption{
		Address:  config.Conf.Common.CommonRedis.RedisAddress,
		Password: config.Conf.Common.CommonRedis.RedisPassword,
		Db:       config.Conf.Common.CommonRedis.Db,
	}
	RedisClient = tools.GetRedisInstance(redisOpt)
	result, err := RedisClient.BRPop(time.Second*10, config.QueueName).Result()
	if err != nil {
		t.Fail()
	}
	t.Log(result, len(result))
	if len(result) >= 1 {
		t.Log(result[1])
	}
}
