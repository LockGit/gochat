/**
 * Created by lock
 * Date: 2019-08-09
 * Time: 15:18
 */
package connect

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gochat/proto"
	"sync"
)

const NoRoom = -1

type Room struct {
	Id          int
	OnlineCount int // room online user count
	rLock       sync.RWMutex
	drop        bool // make room is live
	next        *Channel
}

func NewRoom(roomId int) *Room {
	room := new(Room)
	room.Id = roomId
	room.drop = false
	room.next = nil
	room.OnlineCount = 0
	return room
}

func (r *Room) Put(ch *Channel) (err error) {
	//doubly linked list
	r.rLock.Lock()
	defer r.rLock.Unlock()
	if !r.drop {
		if r.next != nil {
			r.next.Prev = ch
		}
		ch.Next = r.next
		ch.Prev = nil
		r.next = ch
		r.OnlineCount++
	} else {
		err = errors.New("room drop")
	}
	return
}

func (r *Room) Push(msg *proto.Msg) {
	r.rLock.RLock()
	for ch := r.next; ch != nil; ch = ch.Next {
		if err := ch.Push(msg); err != nil {
			logrus.Infof("push msg err:%s", err.Error())
		}
	}
	r.rLock.RUnlock()
	return
}

func (r *Room) DeleteChannel(ch *Channel) bool {
	r.rLock.RLock()
	if ch.Next != nil {
		//if not footer
		ch.Next.Prev = ch.Prev
	}
	if ch.Prev != nil {
		// if not header
		ch.Prev.Next = ch.Next
	} else {
		r.next = ch.Next
	}
	r.OnlineCount--
	r.drop = false
	if r.OnlineCount <= 0 {
		r.drop = true
	}
	r.rLock.RUnlock()
	return r.drop
}
