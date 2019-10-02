/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 15:52
 */
package logic

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gochat/config"
	"gochat/logic/dao"
	"gochat/proto"
	"gochat/tools"
	"strconv"
	"time"
)

type RpcLogic struct {
}

var RedisSessClient *redis.Client

func init() {
	RedisSessClient = NewRedisSessClient()
}

// wrap
func NewRedisSessClient() *redis.Client {
	return RedisClient
}

func (rpc *RpcLogic) Register(ctx context.Context, args *proto.RegisterRequest, reply *proto.RegisterReply) (err error) {
	reply.Code = config.FailReplyCode
	u := new(dao.User)
	u.UserName = args.Name
	userId, err := u.Add()
	if err != nil {
		logrus.Infof("register err:%s", err.Error())
		return
	}
	//set token
	authToken := tools.GetRandomToken("sess", 32)
	err = RedisSessClient.Set(authToken, userId, 86400*time.Second).Err()
	if err != nil {
		logrus.Infof("register set redis token fail!")
		return err
	}
	reply.Code = config.SuccessReplyCode
	reply.AuthToken = authToken
	return
}

func (rpc *RpcLogic) Login(ctx context.Context, args *proto.LoginRequest, reply *proto.LoginResponse) (err error) {
	reply.Code = config.FailReplyCode
	u := new(dao.User)
	userName := args.Name
	data := u.CheckHaveUserName(userName)
	if data.Id == 0 {
		return errors.New("no this user!")
	}
	//set token
	authToken := tools.CreateSessionId()
	err = RedisSessClient.Set(authToken, data.Id, 86400*time.Second).Err()
	if err != nil {
		logrus.Infof("register set redis token fail!")
		return err
	}
	reply.Code = config.SuccessReplyCode
	reply.AuthToken = authToken
	return
}

func (rpc *RpcLogic) CheckAuth(ctx context.Context, args *proto.CheckAuthRequest, reply *proto.CheckAuthResponse) (err error) {
	reply.Code = config.FailReplyCode
	authToken := args.AuthToken
	sessionName := tools.GetSessionName(authToken)
	var userId string
	userId, err = RedisSessClient.Get(sessionName).Result()
	if err != nil {
		logrus.Infof("check auth fail!,token is:%s", authToken)
		return err
	}
	reply.Code = config.SuccessReplyCode
	intUserId, _ := strconv.Atoi(userId)
	reply.UserId = intUserId
	return
}

func (rpc *RpcLogic) Logout(ctx context.Context, args *proto.LogoutRequest, reply *proto.LogoutResponse) (err error) {
	reply.Code = config.FailReplyCode
	authToken := args.AuthToken
	sessionName := tools.GetSessionName(authToken)
	err = RedisSessClient.Del(sessionName).Err()
	if err != nil {
		logrus.Infof("logout error:%s", err.Error())
		return err
	}
	reply.Code = config.SuccessReplyCode
	return
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
