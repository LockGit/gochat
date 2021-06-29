/**
 * Created by lock
 * Date: 2019-08-09
 * Time: 10:56
 */
package main

import (
	"fmt"
	"gochat/connect"
	//"gochat/task"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	connect.New().RunTcp()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Println("Server exiting")
}
