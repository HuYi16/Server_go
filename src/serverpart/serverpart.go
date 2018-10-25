package serverpart

import (
	"encoding/binary"
	"fmt"
	"net"
	"threadpool"
	"unsafe"
)

const (
	GoMode = iota
	PoolMode
	MAX_LEN = 1024
)

type StServerInfo struct {
	Index           int64
	ServerIp        string
	ServerPort      string
	ThreadMode      int //1 use go   2 use threadpool
	DisConnCallBack interface{}
	ReadCallBack    interface{}
	OnlineMap4Self  map[int64]net.Conn
	OnlineMap4Other map[net.Conn]int64
}

var stInfo StServerInfo

func (arg *StServerInfo) AddOnlineInfo(val net.Conn) int64 {
	Index++
	arg.OnlineMap4Self[Index] = val
	arg.OnlineMap4Other[val] = Index
	return Index
}

func (arg *StServerInfo) ReMoveOnlineInfo(conn net.Conn) {
	val, ok := arg.OnlineMap4Other[conn]
	if ok {
		delete(key, arg.OnlineMap4Other)
		val1, ok1 := arg.OnlineMap4Self[val]
		if ok1 {
			delete(val1, arg.OnlineMap4OSelf)
		}
	}
}

func (arg *StServerInfo) GetSocketId(key int64) (net.Conn, bool) {
	return arg.OnlineMap4Self[key]
}

func (arg *StServerInfo) Check() bool {
	if arg.OnlineMap4Self == nil || arg.OnlineMap4Other == nil || ServerIp == "" || ServerPort == "" || DisConCallBack == nil || ReadCallBack == nil {
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

func (arg *StServerInfo) SetCallDisConn(disconn interface{}) bool {
	if disconn != nil {
		arg.DisConnCallBack = disconn
	} else {
		return false
	}
	return true
}

func (arg *StServerInfo) SetCallRead(read interface{}) bool {
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
		OnlineMap4Self:  make(map[int64]net.Conn),
		OnlineMap4Other: make(map[net.Conn]int64),
		Index:           10241111,
		ThreadMode:      GoMode,
	}
}
func SetIpPort(ip, port string) bool {
	return stInfo.SetIpPort(ip, port)
}

func SetCallRead(callback interface{}) bool {
	return stInfo.SetCallRead(callback)
}

func SetCallDisConn(callback interface{}) bool {
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
	info := make([]byte, unsafe.Sizeof(int32))
	len, err := conn.Read(info)
	res := readWriteErr(err)
	if res == 0 {
		stInfo.RemoveOnlineInfo(conn)
		return false
	}
	var msglen int32
	binary.Read(info, binary.BigEndian, &msglen)
	if msglen != 0 {
		tempbuf := make([]byte, msglen+1)
		len, err := conn.Read(tempbuf)
		res := readWriteErr(err)
		if res == 0 {
			stInfo.RemoveOnlineInfo(conn)
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
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", Ip, Port))
	if err != nil {
		fmt.Println("lisetn fail!!", err)
		return false
	}
	for {
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
