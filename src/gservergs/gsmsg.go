package gservergs

import (
	"commondef"
	"fmt"
)

const (
	CLIENT_TYPE = 10000 //client
	PUBLIC_TYPE = 10001
	LOGIC_TYPE  = 10002
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
	return
}

func getMsgHead(head []byte) commondef.MsgHead {
	var info commondef.MsgHead
	return info
}

func ReadData(id int64, msg []byte, len int) {
	fmt.Println("recv data,len:", len)
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

func MsgReGister(msg []byte, len int) {
	return
}
