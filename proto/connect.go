/**
 * Created by lock
 * Date: 2019-08-10
 * Time: 18:00
 */
package proto

type Msg struct {
	Ver       int16  `json:"ver"`  // protocol version
	Operation int32  `json:"op"`   // operation for request
	SeqId     int32  `json:"seq"`  // sequence number chosen by client
	Body      []byte `json:"body"` // binary body bytes
}

type PushMsgRequest struct {
	Uid string
	Msg Msg
}

type PushRoomMsgRequest struct {
	RoomId int
	Msg    Msg
}

type PushRoomCountRequest struct {
	RoomId int
	Count  int
}
