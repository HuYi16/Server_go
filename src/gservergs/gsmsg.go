package gservergs

import (
	//	"commondef"
	"fmt"
	"math/rand"
	"serverpart"
	"time"
	"unsafe"
)

const (
	CLIENT_TYPE = 10000 //client
	PUBLIC_TYPE = 10001
	LOGIC_TYPE  = 10002
	GS_TYPE     = 10003
	REGISTER_ID = 30001
)

type SubHead struct {
	MainId   int
	SubId    int
	SelfType int
}
type NetInfo struct {
	Type        int
	SockId      int64
	LastMsgTime int64
	Key         map[float64]int
}

var OnlineInfo map[int]map[int64]NetInfo

func init() {
	OnlineInfo = make(map[int]map[int64]NetInfo)
}

func checkKey(id int64, head serverpart.MsgHead) bool {
	for _, v := range OnlineInfo {
		vv, ok := v[id]
		if ok {
			_, okkey := vv.Key[head.Key]
			if !okkey {
				DisConn(id)
				return false
			} else {
				return true
			}
		}
	}
	return false
}
func ReadData(id int64, head serverpart.MsgHead, msgBody []byte) {
	fmt.Println("recv data,len:", head.BodyLen)
	if head.ToServerId == GS_TYPE {
		selfMsg(id, msgBody, head.BodyLen)
	} else {
		ok := checkKey(id, head)
		if ok {
			ServerService.Write(id, msgBody, head.ToServerId, head.BodyLen)
			for _, v := range OnlineInfo {
				vv, ok := v[id]
				if ok {
					vv.LastMsgTime = time.Now().Unix()
				}
			}
			return
		} else {
			ServerService.CloseSocket(id)
			return
		}
	}
	return
}

func selfMsg(id int64, msg []byte, len int) {
	head := SubHead{}
	head = **(**SubHead)(unsafe.Pointer(&msg))
	//	head = SubHead(msg[0 : unsafe.Sizeof(head)-1])
	if head.SubId == REGISTER_ID {
		MsgReGister(id, head.SelfType)
	}
	return
}
func DisConn(id int64) {
	fmt.Println("client closed", id)
	for _, v := range OnlineInfo {
		_, ok := v[id]
		if ok {
			delete(v, id)
			return
		}
	}
	return
}

func MsgReGister(id int64, serverId int) {
	//change type and get client type
	ok := false
	//	childMap map[int64]NetInfo
	childMap, ok := OnlineInfo[serverId]
	if !ok {
		childMap = make(map[int64]NetInfo)
	}
	childNode := NetInfo{
		Type:        serverId,
		SockId:      id,
		LastMsgTime: time.Now().Unix(),
	}
	childNode.Key = make(map[float64]int)
	for i := 0; i < 4; i++ {
		childNode.Key[rand.Float64()] = 0
	}
	childMap[id] = childNode
	if !ok {
		OnlineInfo[serverId] = childMap
	}
	return
}
