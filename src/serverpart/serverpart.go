package serverpart

import (
	"fmt"
	"net"
	"sync"
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
	llMaxAcceptNum int64
	llNowAcceptNum int64
	lock           sync.Mutex
	listener       net.Listener
}

func (this *NetCn) addOnlineNode(conn net.Conn) int64 {
	this.index++
	this.onlineMap[this.index] = clientInfo{
		ConnId:      conn,
		LastMsgTime: time.Now().Unix(),
	}
	return this.index
}

func (this *NetCn) nowAcceptNumChange(val int64) {
	this.lock.Lock()
	this.llNowAcceptNum += val
	this.lock.Unlock()
}

func (this *NetCn) removeOnlineNode(conn net.Conn) {
	for k, v := range this.onlineMap {
		if v.ConnId == conn {
			conn.Close()
			if nil != this.disCallBack {
				this.disCallBack(k)
			}
			delete(this.onlineMap, k)
			this.nowAcceptNumChange(-1)
			return
		}
	}
}

func (this *NetCn) getConnId(key int64) (net.Conn, bool) {
	val, ok := this.onlineMap[key]
	if ok {
		return val.ConnId, ok
	}
	return nil, ok
}

func (this *NetCn) checkParam() bool {
	if this.onlineMap == nil || this.ip == "" || this.port == "" || this.disCallBack == nil || this.readCallBack == nil {
		return false
	}
	return true
}

func (this *NetCn) doReadLoop(conn net.Conn, id int64) {
	for {
		res := this.doRead(conn, id)
		if !res {
			break
		}
	}
	this.removeOnlineNode(conn)
	return
}

func (this *NetCn) CloseSocket(index int64) {
	v, ok := this.onlineMap[index]
	if ok {
		this.removeOnlineNode(v.ConnId)
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
		res := this.readWriteErr(err)
		if res == 0 {
			this.removeOnlineNode(conn)
			return nil, false
		}
		nowlen += len
	}
	return info, true
}

func (this *NetCn) updateLastMsg(id int64) {
	if !this.bServerService {
		return
	}
	val, ok := this.onlineMap[id]
	if ok {
		val.LastMsgTime = time.Now().Unix()
	}
}
func (this *NetCn) doRead(conn net.Conn, id int64) bool {
	head := MsgHead{}
	headlen := unsafe.Sizeof(head)
	buf, ok := this.readDataNum(conn, int(headlen))
	if !ok {
		return false
	}
	head = **(**MsgHead)(unsafe.Pointer(&buf))
	this.updateLastMsg(id)
	if head.BodyLen != 0 {
		msgbuf, msgok := this.readDataNum(conn, head.BodyLen)
		if !msgok {
			fmt.Println("read data msg fail!!")
			return false
		} else {
			this.readCallBack(id, head, msgbuf)
			return true
		}
	} else if !this.bServerService && head.BHeartMsg {
		this.sendHeartMsg(conn)
	}
	this.readCallBack(id, head, nil)
	return true
}

func (this *NetCn) checkOnline() {
	for _, v := range this.onlineMap {
		if v.LastMsgTime+10 < time.Now().Unix() {
			this.removeOnlineNode(v.ConnId)
		} else if v.LastMsgTime+5 < time.Now().Unix() {
			this.sendHeartMsg(v.ConnId)
		}
	}
	//use sleep replace ticker  do HeartBeat
	time.Sleep(1 * time.Second)
}
func (this *NetCn) sendHeartMsg(conn net.Conn) {
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
	wmsg := *(*[]byte)(unsafe.Pointer(tempBytes))
	conn.Write(wmsg)
	return
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
	if !this.checkParam() {
		fmt.Println("base ifno not init!!")
		return false
	}
	var err error
	this.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", this.ip, this.port))
	if err != nil {
		fmt.Println("listen fail!!", err)
		return false
	}
	go this.checkOnline()
	for {
		fmt.Println("start accept!!!")
		conn, err := this.listener.Accept()
		if err == nil {
			if this.llMaxAcceptNum == -1 || this.llMaxAcceptNum > this.llNowAcceptNum {
				id := this.addOnlineNode(conn)
				go this.doReadLoop(conn, id)
				this.nowAcceptNumChange(1)
			} else {
				conn.Close()
			}
		} else {
			if this.listener == nil {
				fmt.Println("stop server")
				break
			}
		}
	}
	return true
}

func (this *NetCn) startClient() bool {
	if !this.checkParam() {
		fmt.Println("base ifno not init!!")
		return false
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", this.ip, this.port))
	if err != nil {
		fmt.Println("Dial fail!!!", err)
		return false
	}
	id := this.addOnlineNode(conn)
	go this.doReadLoop(conn, id)
	return true
}

func (this *NetCn) Write(id int64, msg []byte, serverid int, len int) bool {
	val, ok := this.getConnId(id)
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
	wmsg := *(*[]byte)(unsafe.Pointer(tempBytes))
	if len != 0 {
		wmsg = append(wmsg, msg...)
	}
	val.Write(wmsg)
	return true
}

func (this *NetCn) CreateNetCn(ip, port string, disFun DisFunc, readFun ReadFunc, bServerService bool) {
	this.bServerService = bServerService
	this.disCallBack = disFun
	this.readCallBack = readFun
	this.ip = ip
	this.port = port
	this.index = 1024111
	this.onlineMap = make(map[int64]clientInfo)
	this.llMaxAcceptNum = -1 //no limit
	this.llNowAcceptNum = 0
}
func (this *NetCn) Start() bool {
	if this.bServerService {
		return this.startServer()
	}
	return this.startClient()
}

func (this *NetCn) SetMaxAcceptNum(llMaxNum int64) {
	this.llMaxAcceptNum = llMaxNum
}

func (this *NetCn) Stop() {
	for k, v := range this.onlineMap {
		v.ConnId.Close()
		if nil != this.disCallBack {
			this.disCallBack(k)
		}
		delete(this.onlineMap, k)
		this.nowAcceptNumChange(-1)
	}
	this.listener.Close()
	this.listener = nil
}
