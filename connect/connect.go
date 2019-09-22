/**
 * Created by lock
 * Date: 2019-08-09
 * Time: 18:18
 */
package connect

import (
	"github.com/sirupsen/logrus"
	"gochat/config"
	"runtime"
	"time"
)

var DefaultServer *Server

type Connect struct {
}

func New() *Connect {
	return new(Connect)
}

func (c *Connect) Run() {
	// get connect layer config
	connectConfig := config.Conf.Connect

	//set the maximum number of CPUs that can be executing
	runtime.GOMAXPROCS(connectConfig.ConnectBucket.CpuNum)

	//init logic layer rpc client, call logic layer rpc server
	if err := c.InitLogicRpcClient(); err != nil {
		logrus.Panicf("InitLogicRpcClient err:%s", err.Error())
	}
	//init connect layer rpc server, logic client will call this
	Buckets := make([]*Bucket, connectConfig.ConnectBucket.CpuNum)
	for i := 0; i < connectConfig.ConnectBucket.CpuNum; i++ {
		Buckets[i] = NewBucket(BucketOptions{
			ChannelSize:   connectConfig.ConnectBucket.Channel,
			RoomSize:      connectConfig.ConnectBucket.Room,
			RoutineAmount: connectConfig.ConnectBucket.RoutineAmount,
			RoutineSize:   connectConfig.ConnectBucket.RoutineSize,
		})
	}
	operator := new(DefaultOperator)
	DefaultServer = NewServer(Buckets, operator, ServerOptions{
		WriteWait:       10 * time.Second,
		PongWait:        60 * time.Second,
		PingPeriod:      54 * time.Second,
		MaxMessageSize:  512,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		BroadcastSize:   512,
	})

	//init connect layer rpc server ,task layer will call this
	if err := c.InitConnectRpcServer(); err != nil {
		logrus.Panicf("InitConnectRpcServer Fatal error: %s \n", err)
	}

	//start connect layer server handler persistent connection
	if err := c.InitWebsocket(); err != nil {
		logrus.Panicf("connect layer InitWebsocket() error:  %s \n", err.Error())
	}
}
