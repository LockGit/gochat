/**
 * Created by lock
 * Date: 2019-10-06
 * Time: 22:46
 */
package rpc

import (
	"context"
	"gochat/api"
	"gochat/proto"
)

type RpcLogic struct {
}

func (rpc *RpcLogic) Login(req *proto.LoginRequest) (code int, authToken string) {
	reply := proto.LoginResponse{}
	api.LogicRpcClient.Call(context.Background(), "Login", req, reply)
	code = reply.Code
	authToken = reply.AuthToken
	return
}

func (rpc *RpcLogic) Register(req *proto.RegisterRequest) (code int, authToken string) {
	reply := proto.RegisterReply{}
	api.LogicRpcClient.Call(context.Background(), "Register", req, reply)
	code = reply.Code
	authToken = reply.AuthToken
	return
}

func (rpc *RpcLogic) CheckAuth(req *proto.CheckAuthRequest) (code int, userId int) {
	reply := proto.CheckAuthResponse{}
	api.LogicRpcClient.Call(context.Background(), "CheckAuth", req, reply)
	code = reply.Code
	userId = reply.UserId
	return
}

func (rpc *RpcLogic) Logout(req *proto.LogoutRequest) (code int) {
	reply := proto.LogoutResponse{}
	api.LogicRpcClient.Call(context.Background(), "Logout", req, reply)
	code = reply.Code
	return
}

func (rpc *RpcLogic) Push(req *proto.Send) (code int, msg string) {
	reply := proto.SuccessReply{}
	api.LogicRpcClient.Call(context.Background(), "Push", req, reply)
	code = reply.Code
	msg = reply.Msg
	return
}

func (rpc *RpcLogic) PushRoom(req *proto.Send) (code int, msg string) {
	reply := proto.SuccessReply{}
	api.LogicRpcClient.Call(context.Background(), "PushRoom", req, reply)
	code = reply.Code
	msg = reply.Msg
	return
}

func (rpc *RpcLogic) Count(req *proto.Send) (code int, msg string) {
	reply := proto.SuccessReply{}
	api.LogicRpcClient.Call(context.Background(), "Count", req, reply)
	code = reply.Code
	msg = reply.Msg
	return
}

func (rpc *RpcLogic) GetRoomInfo(req *proto.Send) (code int, msg string) {
	reply := proto.SuccessReply{}
	api.LogicRpcClient.Call(context.Background(), "GetRoomInfo", req, reply)
	code = reply.Code
	msg = reply.Msg
	return
}
