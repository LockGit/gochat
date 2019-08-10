/**
 * Created by lock
 * Date: 2019-08-10
 * Time: 18:38
 */
package proto

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
	Op           int    `json:"op"`
}
