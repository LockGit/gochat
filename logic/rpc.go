/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 15:52
 */
package logic

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gochat/config"
	"gochat/proto"
	"strconv"
	"time"
)

type RpcLogic struct {
}

func (rpc *RpcLogic) Connect(ctx context.Context, args *proto.ConnectRequest, reply *proto.ConnectReply) (err error) {
	if args == nil {
		logrus.Errorf("Connect() error(%v)", err)
		return
	}
	logic := new(Logic)
	key := logic.getKey(args.Auth)
	userInfo, err := RedisClient.HGetAll(key).Result()
	if err != nil {
		logrus.Infof("RedisCli HGetAll key :%s , err:%s", key, err)
	}
	reply.Uid = userInfo["UserId"]
	roomUserKey := logic.getRoomUserKey(strconv.Itoa(args.RoomId))
	if reply.Uid == "" {
		reply.Uid = config.NoAuth
	} else {
		userKey := logic.getKey(reply.Uid)
		logrus.Infof("logic redis set userKey:%s, serverId : %s", userKey, args.ServerId)
		validTime := config.RedisBaseValidTime * time.Second
		err = RedisClient.Set(userKey, args.ServerId, validTime).Err()
		if err != nil {
			logrus.Warnf("logic set err:%s", err)
		}
		RedisClient.HSet(roomUserKey, reply.Uid, userInfo["UserName"])
		// add room user count ++
		RedisClient.Incr(logic.getRoomOnlineCountKey(fmt.Sprintf("%d", args.RoomId)))
	}
	logrus.Infof("logic rpc uid:%s", reply.Uid)
	return
}

func (rpc *RpcLogic) DisConnect(ctx context.Context, args *proto.DisConnectRequest, reply *proto.DisConnectReply) (err error) {
	logic := new(Logic)
	roomUserKey := logic.getRoomUserKey(strconv.Itoa(args.RoomId))
	// room user count --
	if args.RoomId > 0 {
		RedisClient.Decr(logic.getRoomOnlineCountKey(fmt.Sprintf("%d", args.RoomId))).Result()
	}
	// room login user--
	if args.Uid != config.NoAuth {
		err = RedisClient.HDel(roomUserKey, args.Uid).Err()
		if err != nil {
			logrus.Warnf("HDel getRoomUserKey err : %s", err)
		}
	}
	//below code can optimize send a signal to queue,another process get a signal from queue,then push event to websocket
	roomUserInfo, err := RedisClient.HGetAll(roomUserKey).Result()
	if err != nil {
		logrus.Warnf("RedisCli HGetAll roomUserInfo key:%s, err: %s", roomUserKey, err)
	}
	if err = logic.RedisPublishRoomInfo(args.RoomId, len(roomUserInfo), roomUserInfo); err != nil {
		logrus.Warnf("publish RedisPublishRoomCount err: %s", err.Error())
		return
	}
	return
}
