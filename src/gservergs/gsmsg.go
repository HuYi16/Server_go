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
	CLIENT_TYPE   = 10000 //client
	PUBLIC_TYPE   = 10001
	LOGIC_TYPE    = 10002
	GS_TYPE       = 10003
	HEART_BEAT_ID = 30000
	REGISTER_ID   = 30001
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
	Key         map[string]int
}

var OnlineInfo map[int]map[int64]NetInfo

func init() {
	OnlineInfo = make(map[int]map[int64]NetInfo)
	go StartCheckOnline()
}

func StartCheckOnline() {
	ticktimer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticktimer.C:
			checkOnline()
			//wait for a end flag
		}
	}
}
func cheackOnline() {
	for k, v := range OnlineInfo {
		for kk, vv := range v {
			if vv.LastMsgTime+5 < time.Now().Unix() {
				SendHeartBeat(kk)
			}
			if vv.LastMsgTime+10 < time.Now().Unix() {
				//send heat beat
				serverpart.CloseSocket(kk)
				delete(v, kk)
				if len(v) == 0 {
					delete(OnlineInfo, k)
				}
				return
			}
		}
	}
	return
}
func SendHeartBeat(id int64) {
	return
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
	fmt.Println("recv data,len:", head.Len)
	if head.ToServerId == GS_TYPE {
		selfMsg(id, msgBody, head.Len)
	} else {
		ok := checkKey(id, head)
		if ok {
			serverpart.Write(id, msgBody, head.ToServerId, head.Len)
			for _, v := range OnlineInfo {
				vv, ok := v[id]
				if ok {
					vv.LastMsgTime = time.Now().Unix()
				}
			}
			return
		} else {
			serverpart.CloseSocket(id)
			return
		}
	}
	return
}

func selfMsg(id int64, msg []byte, len int) {
	head := SubHead{}
	head = **(**SubHead)(unsafe.Pointer(&msg))
	//	head = SubHead(msg[0 : unsafe.Sizeof(head)-1])
	if head.SubId == HEART_BEAT_ID {
		MsgHeartBeatInfo(id, head.SelfType)
	} else if head.SubId == REGISTER_ID {
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

func MsgHeartBeatInfo(id int64, serverid int) {
	v, ok := OnlineInfo[serverid]
	if ok {
		vv, okv := v[id]
		if okv {
			vv.LastMsgTime = time.Now().Unix()
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
	childNode.Key = make(map[string]int)
	for i := 0; i < 4; i++ {
		childNode.Key[fmt.Sprintf("%f", rand.Float64())] = 0
	}
	childMap[id] = childNode
	if !ok {
		OnlineInfo[serverId] = childMap
	}
	return
}
