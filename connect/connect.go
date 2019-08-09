/**
 * Created by lock
 * Date: 2019-08-09
 * Time: 18:18
 */
package connect

import (
	"gochat/config"
	"runtime"
)

type Connect struct {
}

func New() *Connect {
	return new(Connect)
}

func (c *Connect) Run() {
	// read connect layer config
	connectConfig := config.Conf.Connect
	//sets the maximum number of CPUs that can be executing
	runtime.GOMAXPROCS(connectConfig.ConnectBucket.Num)

}
