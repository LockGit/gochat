/**
 * Created by lock
 * Date: 2019-10-06
 * Time: 23:40
 */
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gochat/api/rpc"
	"gochat/proto"
	"gochat/tools"
)

type FormLogin struct {
	UserName string `form:"userName" json:"userName" binding:"required"`
	Password string `form:"passWord" json:"passWord" binding:"required"`
}

func Login(c *gin.Context) {
	var formLogin FormLogin
	if err := c.ShouldBindBodyWith(&formLogin, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	req := &proto.LoginRequest{
		Name:     formLogin.UserName,
		Password: tools.Sha1(formLogin.Password),
	}
	code, authToken, msg := rpc.RpcLogicObj.Login(req)
	if code == tools.CodeFail || authToken == "" {
		tools.FailWithMsg(c, msg)
		return
	}
	tools.SuccessWithMsg(c, "login success", authToken)
}

type FormRegister struct {
	UserName string `form:"userName" json:"userName" binding:"required"`
	Password string `form:"passWord" json:"passWord" binding:"required"`
}

func Register(c *gin.Context) {
	var formRegister FormRegister
	if err := c.ShouldBindBodyWith(&formRegister, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	req := &proto.RegisterRequest{
		Name:     formRegister.UserName,
		Password: tools.Sha1(formRegister.Password),
	}
	code, authToken, msg := rpc.RpcLogicObj.Register(req)
	if code == tools.CodeFail || authToken == "" {
		tools.FailWithMsg(c, msg)
		return
	}
	tools.SuccessWithMsg(c, "register success", authToken)
}

type FormCheckAuth struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func CheckAuth(c *gin.Context) {
	var formCheckAuth FormCheckAuth
	if err := c.ShouldBindBodyWith(&formCheckAuth, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	authToken := formCheckAuth.AuthToken
	req := &proto.CheckAuthRequest{
		AuthToken: authToken,
	}
	code, userId, userName := rpc.RpcLogicObj.CheckAuth(req)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "auth fail")
		return
	}
	var jsonData = map[string]interface{}{
		"userId":   userId,
		"userName": userName,
	}
	tools.SuccessWithMsg(c, "auth success", jsonData)
}

type FormLogout struct {
	AuthToken string `form:"authToken" json:"authToken" binding:"required"`
}

func Logout(c *gin.Context) {
	var formLogout FormLogout
	if err := c.ShouldBindBodyWith(&formLogout, binding.JSON); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	authToken := formLogout.AuthToken
	logoutReq := &proto.LogoutRequest{
		AuthToken: authToken,
	}
	code := rpc.RpcLogicObj.Logout(logoutReq)
	if code == tools.CodeFail {
		tools.FailWithMsg(c, "logout fail!")
		return
	}
	tools.SuccessWithMsg(c, "logout ok!", nil)
}
