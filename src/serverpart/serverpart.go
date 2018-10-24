package serverpart

import (
	"net"
	"fmt"
	"threadpool"
)

const(
	GoMode = iota
	PoolMode
)
type StServerInfo struct{
	OnlineMap4Self chan map[int64]net.Conn
	OnlineMap4Other chan map[net.Conn][int64]
	ServerIp string
	ServerPort string
	ThreadMode int     //1 use go   2 use threadpool
	Id chan int
	DisConnCallBack  interface{}
	ReadCallBack        interface{}
}

func (arg * StServerInfo) Check()bool{
	if arg.ClientOnlineMap == nil || ServerIp == "" || ServerPort == "" || DisConCallBack == nil || ReadCallBack == nil{
		return false
	}
	return true
}
func (arg * StServerInfo) SetIpPort(ip,port string) bool {
	if ip == "" || port == ""{
		return false
	}
	arg.ServerIp = ip
	arg.ServerPort = port
	return true
}

func (arg * StServerInfo) SetCallDisConn(disconn interface{}) bool{
	if disconn != nil{
		arg.DisConn = disconn
	}else{
		return false
	}
	return true
}

func (arg * StServerInfo) SetCallRead(read interface{}) bool{
	if read != nil{
		arg.Read = read
	}else{
		return false
	}
	return true
}

var stInfo StServerInfo

func init() {
	fmt.Println("inti server")
	stInfo = &StServerInfo{
			OnlineMap4Self:make(chan map[int64]net.Conn),
			OnlineMap4Other:make(chan map[net.Conn]int64),
			Id : 10241111.
			ThreadMode : GoMode
		}
}
func SetIpPort(ip,port string) bool{
	return stInfo.SetIpPort(ip,port)
}

func SetCallRead(callback interface{})bool{
	return stInfo.SetCallRead(callback)
}

func SetCallDisConn(callback interface{})bool{
	return stInfo.SetCallDisConn(callback)
}
func SetThreadMode(Mode int) bool{
	if Mode > PoolMode || Mode < GoMode{
		return false
	}
	stInfo.ThreadMode = Mode
	return true
}

func goModeDoReadLoop(conn net.Conn){
	defer conn.Close()
	for{
		res := DoRead(conn)
		if !res{
			break
		}
	}
}
func DoRead(conn net.Conn) bool{
	info := make([]byte,1024)
	len,err := conn.Read(info)
	res := readWriteErr(err)
	if res == 0{
	//	val,ok := stInfo.ClientSocketIdInfo[]
		return 0
	}else if res == -1{
		return false
	}
	if len == 0{
		
	}
	return true
}

func readWriteErr(err error) int{
	if nil != err{
		if err.Error() == "EOF"{
			fmt.Println("client closed")
			return 0
		}else{
			fmt.Println("err ",err.Error())
			return -1
		}
	}
	return 1
}
func StartServer() bool {
	if !stInfo.Check(){
		fmt.Println("server base info not init!!!")
		return false
	} 
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", Ip, Port))
	if err != nil {
		fmt.Println("lisetn fail!!", err)
		return false
	}
	for{
		conn,err := listener.Accept()
		if err == nil{
			if stInfo.ThreadMode == GoMode{
				go DoRead(conn)
			}
		}
	}
	return false
}
