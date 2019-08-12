/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 15:44
 */
package logic

import (
	"bytes"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gochat/config"
	"gochat/proto"
	"gochat/tools"
)

var RedisClient *redis.Client

func (logic *Logic) InitPublish() (err error) {
	redisOpt := tools.RedisOption{
		Host:     "",
		Port:     0,
		Password: "",
		Db:       0,
	}
	RedisClient = tools.GetRedisInstance(redisOpt)
	if pong, err := RedisClient.Ping().Result(); err != nil {
		logrus.Infof("RedisCli Ping Result pong: %s,  err: %s", pong, err)
	}
	return err
}

func (logic *Logic) RedisPublishRoomInfo(roomId int, count int, RoomUserInfo map[string]string) (err error) {
	var redisMsg = &proto.RedisRoomInfo{
		Op:           config.OpRoomInfoSend,
		RoomId:       roomId,
		Count:        count,
		RoomUserInfo: RoomUserInfo,
	}
	redisMsgByte, err := json.Marshal(redisMsg)
	logrus.Infof("RedisPublishRoomInfo redisMsg info : %s", redisMsgByte)
	err = RedisClient.Publish(config.QueueName, redisMsgByte).Err()
	return
}

func (logic *Logic) getRoomUserKey(authKey string) string {
	var returnKey bytes.Buffer
	returnKey.WriteString(config.RedisRoomPrefix)
	returnKey.WriteString(authKey)
	return returnKey.String()
}

func (logic *Logic) getRoomOnlineCountKey(authKey string) string {
	var returnKey bytes.Buffer
	returnKey.WriteString(config.RedisRoomOnlinePrefix)
	returnKey.WriteString(authKey)
	return returnKey.String()
}

func (logic *Logic) getKey(authKey string) string {
	var returnKey bytes.Buffer
	returnKey.WriteString(config.RedisPrefix)
	returnKey.WriteString(authKey)
	return returnKey.String()
}
