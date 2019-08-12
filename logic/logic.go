/**
 * Created by lock
 * Date: 2019-08-09
 * Time: 18:25
 */
package logic

import (
	"gochat/config"
	"log"
	"runtime"
)

type Logic struct {
}

func New() *Logic {
	return new(Logic)
}

func (logic *Logic) Run() {
	//read config
	logicConfig := config.Conf.Logic

	runtime.GOMAXPROCS(logicConfig.LogicBase.CpuNum)

	//init publish redis
	if err := logic.InitPublishRedisClient(); err != nil {
		log.Panicf("logic init publishRedisClient fail,err:%s", err.Error())
	}

	//init rpc server
	if err := logic.InitRpcServer(); err != nil {
		log.Panicf("logic init rpc server fail")
	}

	//init http server

}
