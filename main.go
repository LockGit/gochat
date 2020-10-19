/**
 * Created by lock
 * Date: 2019-08-09
 * Time: 10:56
 */
package main

import (
	"flag"
	"fmt"
	"gochat/api"
	"gochat/connect"
	"gochat/logic"
	"gochat/site"
	"gochat/task"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var module string
	flag.StringVar(&module, "module", "", "assign run module")
	flag.Parse()
	fmt.Println(fmt.Sprintf("start run %s module", module))
	switch module {
	case "logic":
		logic.New().Run()
	case "connect_websocket":
		connect.New().Run()
	case "connect_tcp":
		connect.New().RunTcp()
	case "task":
		task.New().Run()
	case "api":
		api.New().Run()
	case "site":
		site.New().Run()
	default:
		fmt.Println("exiting,module param error!")
		return
	}
	fmt.Println(fmt.Sprintf("run %s module done!", module))
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Println("Server exiting")
}
