package serverpart

import (
	"fmt"
	"net"
	//	"threadpool"
	"strconv"
	"unsafe"
)

const (
	GoMode = iota
	PoolMode
	MAX_LEN = 1024
)

type ReadCallBackFun func(index int64, msg []byte, len int)
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
func DoRead(conn net.Conn, id int64) bool {
	var t int32
	info := make([]byte, unsafe.Sizeof(t))
	_, err := conn.Read(info)
	res := readWriteErr(err)
	if res == 0 {
		stInfo.ReMoveOnlineInfo(conn)
		return false
	}
	msglen, err := strconv.Atoi(string(info))
	if err != nil {
		stInfo.ReMoveOnlineInfo(conn)
		return false
	}
	if msglen != 0 {
		tempbuf := make([]byte, msglen+1)
		_, err := conn.Read(tempbuf)
		res := readWriteErr(err)
		if res == 0 {
			stInfo.ReMoveOnlineInfo(conn)
			return false
		} else {
			stInfo.ReadCallBack(id, tempbuf, msglen)
			return true
		}
	}
	return true
}

func Write(id int64, msg []byte) bool {
	val, ok := stInfo.GetSocketId(id)
	if ok {
		val.Write(msg)
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
		fmt.Println("start Accept")
		conn, err := listener.Accept()
		if err == nil {
			id := stInfo.AddOnlineInfo(conn)
			if stInfo.ThreadMode == GoMode {
				go GoModeDoReadLoop(conn, id)
			} else {
				//thread mode

			}
		}
	}
	return false
}
