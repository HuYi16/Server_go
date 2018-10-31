package serverpart

import (
	"commondef"
	"fmt"
	"net"
	"threadpool"
	"time"
	"unsafe"
)

type MsgHead struct {
	Key        string
	ToServerId int
	Len        int
}

const (
	GoMode = iota
	PoolMode
	MAX_LEN = 1024
)

type ReadCallBackFun func(index int64, head MsgHead, msgBody []byte)
type DisConnCallBackFun func(index int64)

type StServerInfo struct {
	Index           int64
	ServerIp        string
	ServerPort      string
	ThreadMode      int //1 use go   2 use threadpool
	DisConnCallBack DisConnCallBackFun
	ReadCallBack    ReadCallBackFun
	OnlineMap       map[int64]net.Conn
}

var stInfo StServerInfo

func (arg *StServerInfo) AddOnlineInfo(val net.Conn) int64 {
	arg.Index++
	arg.OnlineMap[arg.Index] = val
	return arg.Index
}

func (arg *StServerInfo) ReMoveOnlineInfo(conn net.Conn) {
	for v, k := range arg.OnlineMap {
		if k == conn {
			delete(arg.OnlineMap, v)
			return
		}
	}
}

func (arg *StServerInfo) GetSocketId(key int64) (net.Conn, bool) {
	val, ok := arg.OnlineMap[key]
	return val, ok
}

func (arg *StServerInfo) Check() bool {
	if arg.OnlineMap == nil || arg.ServerIp == "" || arg.ServerPort == "" || arg.DisConnCallBack == nil || arg.ReadCallBack == nil {
		return false
	}
	return true
}
func (arg *StServerInfo) SetIpPort(ip, port string) bool {
	if ip == "" || port == "" {
		return false
	}
	arg.ServerIp = ip
	arg.ServerPort = port
	return true
}

func (arg *StServerInfo) SetCallDisConn(disconn DisConnCallBackFun) bool {
	if disconn != nil {
		arg.DisConnCallBack = disconn
	} else {
		return false
	}
	return true
}

func (arg *StServerInfo) SetCallRead(read ReadCallBackFun) bool {
	if read != nil {
		arg.ReadCallBack = read
	} else {
		return false
	}
	return true
}

func init() {
	fmt.Println("inti server")
	stInfo = StServerInfo{
		OnlineMap:  make(map[int64]net.Conn),
		Index:      10241111,
		ThreadMode: GoMode,
	}
}
func SetIpPort(ip, port string) bool {
	return stInfo.SetIpPort(ip, port)
}

func SetCallRead(callback ReadCallBackFun) bool {
	return stInfo.SetCallRead(callback)
}

func SetCallDisConn(callback DisConnCallBackFun) bool {
	return stInfo.SetCallDisConn(callback)
}
func SetThreadMode(Mode int) bool {
	if Mode > PoolMode || Mode < GoMode {
		return false
	}
	stInfo.ThreadMode = Mode
	return true
}

func GoModeDoReadLoop(conn net.Conn, id int64) {
	defer conn.Close()
	for {
		res := DoRead(conn, id)
		if !res {
			break
		}
	}
	stInfo.DisConnCallBack(id)
	return
}
func ThreadModeReadLoop(arg interface{}) {
	for k, v := range stInfo.OnlineMap {
		v.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
		res := DoRead(v, k)
		if !res {
			v.Close()
		}
	}
}
func CloseSocket(socketid int64) {
	v, ok := stInfo.OnlineMap[socketid]
	if ok {
		delete(stInfo.OnlineMap, socketid)
		v.Close()
	}
	return
}
func readDataNum(conn net.Conn, msglen int) ([]byte, bool) {
	if nil == conn || msglen == 0 {
		return nil, false
	}
	nowlen := 0
	info := make([]byte, msglen)
	for msglen > nowlen {
		len, err := conn.Read(info[nowlen:])
		res := readWriteErr(err)
		if res == 0 {
			stInfo.ReMoveOnlineInfo(conn)
			return nil, false
		}
		nowlen += len
	}
	return info, true
}

type SliceMock struct {
	Addr uintptr
	Len  int
	Cap  int
}

func DoRead(conn net.Conn, id int64) bool {
	head := MsgHead{}
	headlen := unsafe.Sizeof(head)
	buf, ok := readDataNum(conn, int(headlen))
	if !ok {
		return false
	}
	/*
		tempBytes := &SliceMock{
			Addr: uintptr(unsafe.Pointer(&head)),
			Len:  int(headlen),
			Cap:  int(headlen),
		}
		data := *(*[]byte)(unsafe.Pointer(tempBytes))
	*/
	head = **(**MsgHead)(unsafe.Pointer(&buf))
	if head.Len != 0 {
		msgbuf, msgok := readDataNum(conn, head.Len)
		if !msgok {
			fmt.Println("read data msg fail!!")
			return false
		} else {
			stInfo.ReadCallBack(id, head, msgbuf)
			return true
		}
	}
	return true
}

func Write(id int64, msg []byte, serverid int, len int) bool {
	val, ok := stInfo.GetSocketId(id)
	head := MsgHead{
		Len:        len,
		ToServerId: serverid,
	}
	headlen := unsafe.Sizeof(head)
	tempBytes := &SliceMock{
		Addr: uintptr(unsafe.Pointer(&head)),
		Len:  int(headlen),
		Cap:  int(headlen),
	}
	wmsg := *(*[]byte)(unsafe.Pointer(tempBytes))
	wmsg = append(wmsg, msg...)
	if ok {
		val.Write(wmsg)
		return true
	}
	return false

}
func readWriteErr(err error) int {
	if nil != err {
		if err.Error() == "EOF" {
			fmt.Println("client closed")
			return 0
		} else {
			fmt.Println("err ", err.Error())
			return -1
		}
	}
	return 1
}
func StartServer() bool {
	if !stInfo.Check() {
		fmt.Println("server base info not init!!!")
		return false
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", stInfo.ServerIp, stInfo.ServerPort))
	if err != nil {
		fmt.Println("lisetn fail!!", err)
		return false
	}
	for {
		fmt.Println("start accept")
		conn, err := listener.Accept()
		if err == nil {
			id := stInfo.AddOnlineInfo(conn)
			if stInfo.ThreadMode == GoMode {
				go GoModeDoReadLoop(conn, id)
				return true
			} else {
				//thread mode
				job := commondef.StJobInfo{
					RepeatTimes: -1,
					Job:         ThreadModeReadLoop,
				}
				threadpool.AddTask(job)
				return true
			}
		}
	}
	return false
}
