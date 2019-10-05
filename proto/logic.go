/**
 * Created by lock
 * Date: 2019-08-10
 * Time: 18:38
 */
package proto

type LoginRequest struct {
	Name string
}

type LoginResponse struct {
	Code      int
	AuthToken string
}

type RegisterRequest struct {
	Name string
}

type RegisterReply struct {
	Code      int
	AuthToken string
}

type LogoutRequest struct {
	AuthToken string
}

type LogoutResponse struct {
	Code int
}

type CheckAuthRequest struct {
	AuthToken string
}

type CheckAuthResponse struct {
	Code   int
	UserId int
}

type ConnectRequest struct {
	Auth     string
	RoomId   int
	ServerId int
}

type ConnectReply struct {
	Uid string
}

type DisConnectRequest struct {
	RoomId int
	Uid    string
}

type DisConnectReply struct {
	Has bool
}

type Send struct {
	Code         int    `json:"code"`
	Msg          string `json:"msg"`
	FormUserId   string `json:"formUserId"`
	FormUserName string `json:"formUserName"`
	ToUserId     string `json:"toUserId"`
	ToUserName   string `json:"toUserName"`
	RoomId       int    `json:"roomId"`
	Op           int    `json:"op"`
}
