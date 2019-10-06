/**
 * Created by lock
 * Date: 2019-10-06
 * Time: 23:40
 */
package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gochat/api/rpc"
	"gochat/proto"
)

func Login(c *gin.Context) {
	rpcLogic := new(rpc.RpcLogic)
	req := &proto.LoginRequest{
		Name: "test",
	}
	code, authToken := rpcLogic.Login(req)
	fmt.Println(code, authToken)
}

func Register(c *gin.Context) {

}

func CheckAuth(c *gin.Context) {

}

func Logout(c *gin.Context) {

}
