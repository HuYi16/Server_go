package gservergs

import (
	"commondef"
	"fmt"
	"serverpart"
	"time"
)

const (
	CLIENT_TYPE = 10000 //client
	PUBLIC_TYPE = 10001
	LOGIC_TYPE  = 10002
	GS_TYPE     = 10003
)

type NetInfo struct {
	Type        int
	SockId      int64
	LastMsgTime int64
}

var OnlineInfo map[int]map[int64]NetInfo

func init() {
	OnlineInfo = make(map[int]map[int64]NetInfo)
}

func cheackOnline() {
	for _, v := range OnlineInfo {
		for _, vv := range v {
			if vv.LastMsgTime+10 < time.Now().Unix() {
				//send heat beat
				CloseSicket(vv)
				delete(v, vv)
			}
		}
	}
	return
}

func getMsgHead(head []byte) commondef.MsgHead {
	var info commondef.MsgHead
	info = commondef.MsgHead(head)
	return info
}

func ReadData(id int64, msg []byte, len int) {
	fmt.Println("recv data,len:", len)
	head := getMsgHead(msg)
	if head.ToServerId == GS_TYPE {
		MsgReGister(msg, len, id)
	}
	return
}

func DisConn(id int64) {
	fmt.Println("client closed", id)
	return
}

func dispatchMsg(msg []byte, len int) {
	return
}

func MsgHeartBeatInfo(msg []byte, len int) {
	return
}

func MsgOther(msg []byte, len int) {
	return
}

func MsgReGister(msg []byte, len int, id int64) {
	//change type and get client type
	clientType := 1
	ok := false
	var childMap map[int64]NetInfo
	childMap, ok = OnlineInfo[clinetType]
	if !ok {
		childMap := make(map[int64]NetInfo)
	}
	childMap[id] = NetInfo{
		Type:        1,
		SockId:      id,
		LastMsgTime: time.Now().Unix(),
	}
	if !ok {
		OnlineMap[clientType] = childMap
	}
	return
}
