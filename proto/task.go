/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 19:23
 */
package proto

type RedisMsg struct {
	Op           int               `json:"op"`
	ServerId     int               `json:"serverId,omitempty"`
	RoomId       int               `json:"roomId,omitempty"`
	UserId       string            `json:"userId,omitempty"`
	Msg          []byte            `json:"msg"`
	Count        int               `json:"count"`
	RoomUserInfo map[string]string `json:"RoomUserInfo"`
}

type RedisRoomInfo struct {
	Op           int               `json:"op"`
	RoomId       int               `json:"roomId,omitempty"`
	Count        int               `json:"count,omitempty"`
	RoomUserInfo map[string]string `json:"roomUserInfo"`
}

type SuccessReply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
