/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 15:52
 */
package logic

import (
	"context"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"gochat/config"
	"gochat/proto"
	"gochat/tools"
	"strconv"
	"strings"
	"time"
)

type LogicRpc struct {
}

func (logic *Logic) InitRpcServer() (err error) {
	var network, addr string
	// a host multi port case
	rpcAddressList := strings.Split(config.Conf.Logic.LogicBase.RpcAddress, ",")
	for _, bind := range rpcAddressList {
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Panicf("InitLogicRpc ParseNetwork error : %s", err.Error())
		}
		go logic.createRpcServer(network, addr)
	}
	return
}

func (logic *Logic) createRpcServer(network string, addr string) {
	s := server.NewServer()
	logic.addRegistryPlugin(s, network, addr)
	// serverId must be unique
	s.RegisterName(config.Conf.Common.CommonEtcd.ServerPathLogic, new(LogicRpc), fmt.Sprintf("%d", config.Conf.Common.CommonEtcd.ServerId))
	s.Serve(network, addr)
}

func (logic *Logic) addRegistryPlugin(s *server.Server, network string, addr string) {
	r := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: network + "@" + addr,
		EtcdServers:    []string{config.Conf.Common.CommonEtcd.Host},
		BasePath:       config.Conf.Common.CommonEtcd.BasePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		logrus.Fatal(err)
	}
	s.Plugins.Add(r)
}

func (rpc *LogicRpc) Connect(ctx context.Context, args *proto.ConnectRequest, reply *proto.ConnectReply) (err error) {
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

func (rpc *LogicRpc) Disconnect(ctx context.Context, args *proto.DisConnectRequest, reply *proto.DisConnectReply) (err error) {
	logic := new(Logic)
	roomUserKey := logic.getRoomUserKey(strconv.Itoa(int(args.RoomId)))

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
	roomUserInfo, err := RedisClient.HGetAll(roomUserKey).Result()
	if err != nil {
		logrus.Warnf("RedisCli HGetAll roomUserInfo key:%s, err: %s", roomUserKey, err)
	}
	if err = logic.RedisPublishRoomInfo(args.RoomId, len(roomUserInfo), roomUserInfo); err != nil {
		logrus.Warnf("Count redis RedisPublishRoomCount err: %s", err)
		return
	}
	return
}
