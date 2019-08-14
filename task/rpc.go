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

var (
	logicRpcClient client.XClient
	RpcClientList  map[int]client.XClient
)

func (task *Task) InitConnectRpcClient() (err error) {
	etcdConfig := config.Conf.Common.CommonEtcd
	d := client.NewEtcdDiscovery(etcdConfig.BasePath, etcdConfig.ServerPathComet, []string{Conf.ZookeeperInfo.Host}, nil)

	RpcClientList = make(map[int]client.XClient, len(d.GetServices()))
	// Get comet service configuration from zookeeper
	for _, cometConf := range d.GetServices() {

		cometConf.Value = strings.Replace(cometConf.Value, "=&tps=0", "", 1)

		serverId, error := strconv.ParseInt(cometConf.Value, 10, 8)
		if error != nil {
			logrus.Panicf("InitComets err，Can't find serverId. error: %s", error)
		}
		d := client.NewPeer2PeerDiscovery(cometConf.Key, "")
		RpcClientList[int8(serverId)] = client.NewXClient(Conf.ZookeeperInfo.ServerPathComet, client.Failtry, client.RandomSelect, d, client.DefaultOption)
		log.Infof("RpcClientList addr %s, v %v", cometConf.Key, RpcClientList[int8(serverId)])

	}
	logicRpcClient = client.NewXClient(Conf.ZookeeperInfo.ServerPathComet, client.Failtry, client.RandomSelect, d, client.DefaultOption)

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
