/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 23:36
 */
package connect

import (
	"context"
	"errors"
	"fmt"
	"github.com/rcrowley/go-metrics"
	"github.com/rpcxio/libkv/store"
	etcdV3 "github.com/rpcxio/rpcx-etcd/client"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/server"
	"gochat/config"
	"gochat/proto"
	"gochat/tools"
	"strings"
	"sync"
	"time"
)

var logicRpcClient client.XClient
var once sync.Once

type RpcConnect struct {
}

func (c *Connect) InitLogicRpcClient() (err error) {
	etcdConfigOption := &store.Config{
		ClientTLS:         nil,
		TLS:               nil,
		ConnectionTimeout: time.Duration(config.Conf.Common.CommonEtcd.ConnectionTimeout) * time.Second,
		Bucket:            "",
		PersistConnection: true,
		Username:          config.Conf.Common.CommonEtcd.UserName,
		Password:          config.Conf.Common.CommonEtcd.Password,
	}
	once.Do(func() {
		d, e := etcdV3.NewEtcdV3Discovery(
			config.Conf.Common.CommonEtcd.BasePath,
			config.Conf.Common.CommonEtcd.ServerPathLogic,
			[]string{config.Conf.Common.CommonEtcd.Host},
			true,
			etcdConfigOption,
		)
		if e != nil {
			logrus.Fatalf("init connect rpc etcd discovery client fail:%s", e.Error())
		}
		logicRpcClient = client.NewXClient(config.Conf.Common.CommonEtcd.ServerPathLogic, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	})
	if logicRpcClient == nil {
		return errors.New("get rpc client nil")
	}
	return
}

func (rpc *RpcConnect) Connect(connReq *proto.ConnectRequest) (uid int, err error) {
	reply := &proto.ConnectReply{}
	err = logicRpcClient.Call(context.Background(), "Connect", connReq, reply)
	if err != nil {
		logrus.Fatalf("failed to call: %v", err)
	}
	uid = reply.UserId
	logrus.Infof("connect logic userId :%d", reply.UserId)
	return
}

func (rpc *RpcConnect) DisConnect(disConnReq *proto.DisConnectRequest) (err error) {
	reply := &proto.DisConnectReply{}
	if err = logicRpcClient.Call(context.Background(), "DisConnect", disConnReq, reply); err != nil {
		logrus.Fatalf("failed to call: %v", err)
	}
	return
}

func (c *Connect) InitConnectWebsocketRpcServer() (err error) {
	var network, addr string
	connectRpcAddress := strings.Split(config.Conf.Connect.ConnectRpcAddressWebSockts.Address, ",")
	for _, bind := range connectRpcAddress {
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Panicf("InitConnectWebsocketRpcServer ParseNetwork error : %s", err)
		}
		logrus.Infof("Connect start run at-->%s:%s", network, addr)
		go c.createConnectWebsocktsRpcServer(network, addr)
	}
	return
}

func (c *Connect) InitConnectTcpRpcServer() (err error) {
	var network, addr string
	connectRpcAddress := strings.Split(config.Conf.Connect.ConnectRpcAddressTcp.Address, ",")
	for _, bind := range connectRpcAddress {
		if network, addr, err = tools.ParseNetwork(bind); err != nil {
			logrus.Panicf("InitConnectTcpRpcServer ParseNetwork error : %s", err)
		}
		logrus.Infof("Connect start run at-->%s:%s", network, addr)
		go c.createConnectTcpRpcServer(network, addr)
	}
	return
}

type RpcConnectPush struct {
}

func (rpc *RpcConnectPush) PushSingleMsg(ctx context.Context, pushMsgReq *proto.PushMsgRequest, successReply *proto.SuccessReply) (err error) {
	var (
		bucket  *Bucket
		channel *Channel
	)
	logrus.Info("rpc PushMsg :%v ", pushMsgReq)
	if pushMsgReq == nil {
		logrus.Errorf("rpc PushSingleMsg() args:(%v)", pushMsgReq)
		return
	}
	bucket = DefaultServer.Bucket(pushMsgReq.UserId)
	if channel = bucket.Channel(pushMsgReq.UserId); channel != nil {
		err = channel.Push(&pushMsgReq.Msg)
		logrus.Infof("DefaultServer Channel err nil ,args: %v", pushMsgReq)
		return
	}
	successReply.Code = config.SuccessReplyCode
	successReply.Msg = config.SuccessReplyMsg
	logrus.Infof("successReply:%v", successReply)
	return
}

func (rpc *RpcConnectPush) PushRoomMsg(ctx context.Context, pushRoomMsgReq *proto.PushRoomMsgRequest, successReply *proto.SuccessReply) (err error) {
	successReply.Code = config.SuccessReplyCode
	successReply.Msg = config.SuccessReplyMsg
	logrus.Infof("PushRoomMsg msg %+v", pushRoomMsgReq)
	for _, bucket := range DefaultServer.Buckets {
		bucket.BroadcastRoom(pushRoomMsgReq)
	}
	return
}

func (rpc *RpcConnectPush) PushRoomCount(ctx context.Context, pushRoomMsgReq *proto.PushRoomMsgRequest, successReply *proto.SuccessReply) (err error) {
	successReply.Code = config.SuccessReplyCode
	successReply.Msg = config.SuccessReplyMsg
	logrus.Infof("PushRoomCount msg %v", pushRoomMsgReq)
	for _, bucket := range DefaultServer.Buckets {
		bucket.BroadcastRoom(pushRoomMsgReq)
	}
	return
}

func (rpc *RpcConnectPush) PushRoomInfo(ctx context.Context, pushRoomMsgReq *proto.PushRoomMsgRequest, successReply *proto.SuccessReply) (err error) {
	successReply.Code = config.SuccessReplyCode
	successReply.Msg = config.SuccessReplyMsg
	logrus.Infof("connect,PushRoomInfo msg %+v", pushRoomMsgReq)
	for _, bucket := range DefaultServer.Buckets {
		bucket.BroadcastRoom(pushRoomMsgReq)
	}
	return
}

func (c *Connect) createConnectWebsocktsRpcServer(network string, addr string) {
	s := server.NewServer()
	addRegistryPlugin(s, network, addr)
	//config.Conf.Connect.ConnectTcp.ServerId
	//s.RegisterName(config.Conf.Common.CommonEtcd.ServerPathConnect, new(RpcConnectPush), fmt.Sprintf("%s", config.Conf.Connect.ConnectWebsocket.ServerId))
	s.RegisterName(config.Conf.Common.CommonEtcd.ServerPathConnect, new(RpcConnectPush), fmt.Sprintf("serverId=%s&serverType=ws", c.ServerId))
	s.RegisterOnShutdown(func(s *server.Server) {
		s.UnregisterAll()
	})
	s.Serve(network, addr)
}

func (c *Connect) createConnectTcpRpcServer(network string, addr string) {
	s := server.NewServer()
	addRegistryPlugin(s, network, addr)
	//s.RegisterName(config.Conf.Common.CommonEtcd.ServerPathConnect, new(RpcConnectPush), fmt.Sprintf("%s", config.Conf.Connect.ConnectTcp.ServerId))
	s.RegisterName(config.Conf.Common.CommonEtcd.ServerPathConnect, new(RpcConnectPush), fmt.Sprintf("serverId=%s&serverType=tcp", c.ServerId))
	s.RegisterOnShutdown(func(s *server.Server) {
		s.UnregisterAll()
	})
	s.Serve(network, addr)
}

func addRegistryPlugin(s *server.Server, network string, addr string) {
	r := &serverplugin.EtcdV3RegisterPlugin{
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
