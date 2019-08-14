/**
 * Created by lock
 * Date: 2019-08-13
 * Time: 10:13
 */
package task

import (
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/client"
	"gochat/config"
	"strconv"
	"strings"
)

var RpcConnectClientList map[int]client.XClient

func (task *Task) InitConnectRpcClient() (err error) {
	etcdConfig := config.Conf.Common.CommonEtcd
	d := client.NewEtcdDiscovery(etcdConfig.BasePath, etcdConfig.ServerPathConnect, []string{etcdConfig.Host}, nil)
	if len(d.GetServices()) <= 0 {
		logrus.Panicf("no etcd server find!")
	}
	RpcConnectClientList = make(map[int]client.XClient, len(d.GetServices()))
	for _, cometConf := range d.GetServices() {
		cometConf.Value = strings.Replace(cometConf.Value, "=&tps=0", "", 1)
		serverId, error := strconv.ParseInt(cometConf.Value, 10, 8)
		if error != nil {
			logrus.Panicf("InitComets err，Can't find serverId. error: %s", error)
		}
		d := client.NewPeer2PeerDiscovery(cometConf.Key, "")
		RpcConnectClientList[int(serverId)] = client.NewXClient(etcdConfig.ServerPathConnect, client.Failtry, client.RandomSelect, d, client.DefaultOption)
		logrus.Infof("InitConnectRpcClient addr %s, v %v", cometConf.Key, RpcConnectClientList[int(serverId)])
	}
	return
}

func (task *Task) pushSingleToConnect(serverId int, userId string, msg []byte) {

}

func (task *Task) broadcastRoomToConnect(roomId int, msg []byte) {

}

func (task *Task) broadcastRoomCountToConnect(roomId, count int) {

}

func (task *Task) broadcastRoomInfoToConnect(roomId int, roomUserInfo map[string]string) {

}
