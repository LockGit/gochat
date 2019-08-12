/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 23:36
 */
package connect

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"gochat/config"
	"gochat/proto"
	"sync"
)

var logicRpcClient client.XClient
var once sync.Once

type RpcConnect struct {
}

func InitLogicRpcClient() (err error) {
	once.Do(func() {
		d := client.NewEtcdDiscovery(
			config.Conf.Common.CommonEtcd.BasePath,
			config.Conf.Common.CommonEtcd.ServerPathLogic,
			[]string{config.Conf.Common.CommonEtcd.Host},
			nil,
		)
		logicRpcClient = client.NewXClient(config.Conf.Common.CommonEtcd.ServerPathLogic, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	})
	if logicRpcClient == nil {
		return errors.New("get rpc client nil")
	}
	return
}

func (rpc *RpcConnect) connect(connReq *proto.ConnectRequest) (uid string, err error) {
	reply := &proto.ConnectReply{}
	err = logicRpcClient.Call(context.Background(), "Connect", connReq, reply)
	if err != nil {
		logrus.Fatalf("failed to call: %v", err)
	}
	uid = reply.Uid
	logrus.Infof("comet logic uid :%s", reply.Uid)
	return
}

func (rpc *RpcConnect) disConnect(disConnReq *proto.DisConnectRequest) (err error) {
	reply := &proto.DisConnectReply{}
	if err = logicRpcClient.Call(context.Background(), "DisConnect", disConnReq, reply); err != nil {
		logrus.Fatalf("failed to call: %v", err)
	}
	return
}
