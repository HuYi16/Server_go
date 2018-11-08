package serverpart

import (
	"fmt"
	"net"
	"time"
	"unsafe"
)

type FakeSlice struct {
	addr uintptr
	len  int
	cap  int
}

type MsgHead struct {
	Key        float64
	ToServerId int
	BodyLen    int
	BHeartMsg  bool
}

type clientInfo struct {
	ConnId      net.Conn
	LastMsgTime int64
}

const (
	MAX_LEN = 1024
)

type ReadFunc func(index int64, head MsgHead, msgBody []byte)
type DisFunc func(index int64)

type NetCn struct {
	index          int64
	bServerService bool
	ip             string
	port           string
	disCallBack    DisFunc
	readCallBack   ReadFunc
	onlineMap      map[int64]clientInfo
}

func (this *NetCn) CreateNetCn(ip, port string, disFun DisFunc, readfun ReadFunc, bServerService bool) {
	this.bServerService = bServerService
	this.disCallBack = disFun
	this.readCallBack = readFunc
	this.ip = ip
	this.port = port
	this.index = 1024111
	onlineMap = make(map[int64]clientInfo)
}

func (this *NetCn) addOnlineNode(conn net.Conn) int64 {
	this.index++
	this.onLineMap[this.index] = clinetInfo{
		ConnId:      conn,
		LastMsgTime: time.Now().Unix(),
	}
	return this.index
}

func (this *NetCn) removeOnlineNode(conn net.Conn) {
	for k, v := range this.onlineMap[key] {
		if v.ConnId == conn {
			conn.Close()
			if nil != this.disCallBack {
				this.disCallBack(k)
			}
			delete(this.onlineMap, k)
			return
		}
	}
}

func (this *NetCn) getConnId(key int64) (net.Conn, bool) {
	val, ok := this.onlineMap[key]
	if ok {
		return val.ConnId, ok
	}
	return val, ok
}

func (this *NetCn) checkParam() bool {
	if this.onlineMap == nil || this.Ip == "" || this.Port == "" || DisCallBack == nil || ReadCallBack == nil {
		return false
	}
	return true
}

func (this *NetCn) doReadLoop(conn net.Conn, id int64) {
	for {
		res := doRead(conn, id)
		if !res {
			break
		}
	}
	return
}

func (this *NetCn) closeSocket(index int64) {
	v, ok := this.onlineMap[index]
	if ok {
		this.removeOnlineInfo(v.ConnId)
		return
	}
	return
}

func (this *NetCn) readDataNum(conn net.Conn, msglen int) ([]byte, bool) {
	if nil == conn || 0 == msglen {
		return nil, false
	}
	nowlen := 0
	info := make([]byte, msglen)
	for msglen > nowlen {
		len, err := conn.Read(info[nowlen:])
		res := readWriteErr(err)
		if res == 0 {
			this.removeOnlineInfo(conn)
			return nil, false
		}
		nowlen += len
	}
	return info, true
}

func (this *NetCn) updateLastMsg(id) {
	if !this.bServerService {
		return
	}
	val, ok := this.onLineMap[id]
	if ok {
		val.LastMsgTime = time.Now().Unix()
	}
}
func (this *NetCn) doRead(conn net.Conn, id int64) bool {
	head := MsgHead{}
	headlen := unsafe.Sizeof(head)
	buf, ok := readdataNum(conn, int(headlen))
	if !ok {
		return false
	}
	head = **(**msgHead)(unsafe.Pointer(&buf))
	this.updateLastMsg(id)
	if head.BodyLen != 0 {
		msgbuf, msgok := readDataNum(conn, head.BodyLen)
		if !msgok {
			fmt.Println("read data msg fail!!")
			return false
		} else {
			this.readCallBack(id, head, msgbuf)
			return true
		}
	} else if !this.bServerService && this.BHeartMsg {
		this.SendHeartMsg(conn)
	}
	this.readCallback(id, head, nil)
	return true
}

func (this *NetCn) checkOnline() {
	for k, v := range this.onlineMap {
		if v.LastMsgTime+10 < time.Now().Unix() {
			this.removeOnlienInfo(v.ConnId)
		} else if v.LastMsgTime+5 < time.Now().Unix() {
			this.SendHeartMsg(v.ConnId)
		}
	}
	//use sleep replace ticker  do HeartBeat
	time.Sleep(1 * time.Second)
}
func (this *NetCn) SendHeartMsg(conn net.Conn) {
	head := MsgHead{
		BodyLen:    0,
		ToServerId: 1,
		BHeartMsg:  true,
	}
	headlen := unsafe.Sizeof(head)
	tempBytes := &FakeSlice{
		addr: uintptr(unsafe.Pointer(&head)),
		len:  int(headlen),
		cap:  int(headlen),
	}
	wmsg := *(*[]byte)(unsafe.Pointer(tempByte))
	val.Write(wmsg)
	return true
}
func (this *NetCn) Write(id int64, msg []byte, serverid int, len int) bool {
	val, ok := this.getSocketId(id)
	if !ok {
		return false
	}
	head := MsgHead{
		BodyLen:    len,
		ToServerId: serverid,
		BHeartMsg:  false,
	}
	headlen := unsafe.Sizeof(head)
	tempBytes := &FakeSlice{
		addr: uintptr(unsafe.Pointer(&head)),
		len:  int(headlen),
		cap:  int(headlen),
	}
	wmsg := *(*[]byte)(unsafe.Pointer(tempByte))
	if len != 0 {
		wmgs := append(wmsg, msg...)
	}
	val.Write(wmsg)
	return true
}

func (this *NetCn) readWriteErr(err error) int {
	if nil != err {
		if err.Error() == "EOF" {
			fmt.Println("client closed!!!")
			return 0
		} else {
			return -1
		}
	}
	return 1
}

func (this *NetCn) startServer() bool {
	if !this.check() {
		fmt.Println("base ifno not init!!")
		return false
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", this.ip, this.port))
	if err != nil {
		fmt.Println("listen fail!!", err)
		return false
	}
	go checkOnline()
	for {
		fmt.Println("start accept!!!")
		conn, err := listener.Accept()
		if err == nil {
			id := this.addOnlineNode(conn)
			go doReadLoop(conn, id)
		}
	}
	return true
}

func (this *NetCn) startClient() bool {
	if !this.check() {
		fmt.Println("base ifno not init!!")
		return false
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", this.ip, this.port))
	if err != nil {
		fmt.Println("Dial fail!!!", err)
		return false
	}
	id := this.addOnlineNode(conn)
	go doReadLoop(conn, id)
	return true
}

func (this *NetCn) Start() bool {
	if this.bServerService {
		return startServer()
	}
	return startClient()
}
