/**
 * Created by lock
 * Date: 2019-08-09
 * Time: 15:18
 */
package connect

import (
	"gochat/proto"
	"sync"
)

type Bucket struct {
	cLock         sync.RWMutex        // protect the channels for chs
	chs           map[string]*Channel // map sub key to a channel
	bucketOptions BucketOptions
	rooms         map[int]*Room // bucket room channels
	routines      []chan *proto.PushRoomMsgRequest
	routinesNum   uint64
	broadcast     chan []byte
}

type BucketOptions struct {
	ChannelSize   int
	RoomSize      int
	RoutineAmount uint64
	RoutineSize   int
}

func NewBucket(bucketOptions BucketOptions) (b *Bucket) {
	b = new(Bucket)
	b.chs = make(map[string]*Channel, bucketOptions.ChannelSize)
	b.bucketOptions = bucketOptions
	b.routines = make([]chan *proto.PushRoomMsgRequest, bucketOptions.RoutineAmount)
	b.rooms = make(map[int]*Room, bucketOptions.RoomSize)
	for i := uint64(0); i < b.bucketOptions.RoutineAmount; i++ {
		c := make(chan *proto.PushRoomMsgRequest, bucketOptions.RoutineSize)
		b.routines[i] = c
		go b.PushRoom(c)
	}
	return
}

func (b *Bucket) Put(uid string, roomId int, ch *Channel) (err error) {
	var (
		room *Room
		ok   bool
	)
	b.cLock.Lock()
	if roomId != NoRoom {
		if room, ok = b.rooms[roomId]; !ok {
			room = NewRoom(roomId)
			b.rooms[roomId] = room
		}
		ch.Room = room
	}
	ch.uid = uid
	b.chs[uid] = ch
	b.cLock.Unlock()

	if room != nil {
		err = room.Put(ch)
	}
	return
}

func (b *Bucket) DeleteChannel(ch *Channel) {
	var (
		ok   bool
		room *Room
	)
	b.cLock.RLock()
	if ch, ok = b.chs[ch.uid]; ok {
		room = b.chs[ch.uid].Room
		//delete from bucket
		delete(b.chs, ch.uid)
	}
	if room != nil && room.DeleteChannel(ch) {
		// if room empty delete,will mark room.drop is true
		if room.drop == true {
			delete(b.rooms, room.Id)
		}
	}
	b.cLock.RUnlock()
}

func (b *Bucket) PushRoom(c chan *proto.PushRoomMsgRequest) {

}
