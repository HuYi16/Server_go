package gservergs

import (
	"commondef"
	"fmt"
	//	"net"
	"redispack"
	"serverpart"
	"sqlpart"
	"threadpool"
	//	"time"
)

var ServerInfo commondef.StServerInfo
var ClientSocketInfoMap chan map[int]int64
var ClientUserIdInfoMap chan map[int]commondef.ClientNetInfo

func LoadConfig() {
	ServerInfo.Ip = "47.106.141.213"
	ServerInfo.Port = "8001"
	ServerInfo.NowNumber = 0
	ServerInfo.BalanceNumber = 100
	ServerInfo.BtempLock = false
}
func init() {
	LoadConfig()
	ClientSocketInfoMap = make(chan map[int]int64)
	ClientUserIdInfoMap = make(chan map[int]commondef.ClientNetInfo)
}

/*
type test struct {
	Id   int
	Name string
}

func Job1(arg interface{}) {
	val, ok := arg.(test)
	if !ok {
		fmt.Println("mode 1 type err")
		return
	}
	fmt.Println("mode 1 arg is ", val.Id, "--", val.Name)
	switch val := arg.(type) {
	case test:
		fmt.Println("mode 2 arg is", val.Id, val.Name)
	default:
		fmt.Println("mode 2 arg err")
	}
	time.Sleep(1 * time.Second)
}
func Job2(arg interface{}) {
	fmt.Println("job2 test")
	time.Sleep(3 * time.Second)
}
func Job3(arg interface{}) {
	fmt.Println("job3 test")
	time.Sleep(2 * time.Second)
}
*/
/*
func JobReadData() {

}
func StartServer() bool {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ServerInfo.Ip, ServerInfo.Port))
	if err != nil {
		fmt.Println("listen fail!!")
		panic(err)
		return false
	}
	job := commondef.StJobInfo{
		RepeatTimes: -1,
		Job:         JobReadData,
	}
	for {
		conn, err := listener.Accept()
		if nil != err {
			fmt.Println("accept client fail", err)
		} else {

		}
	}
}
*/
func ReadData(id int64, msg []byte, len int) {
	fmt.Println("recv data_len:", len)
}

func DisConn(id int64) {
	fmt.Println("client closed", id)
}

func StartGs() bool {
	go StartTimer()
	//redis test
	redispack.SetRedisBaseInfo("127.0.0.1:6379", "", "")
	redispack.RedisSet(1, "test", "suc")
	redispack.RedisSet(1, "test1", "succ")
	ok, val := redispack.RedisGet(1, "test")
	fmt.Println("test", ok, val)
	//	threadpool.SetThreadNum(3)
	threadpool.StartThreadPool()
	key := sqlpart.StartSql(commondef.StSqlRedisBaseInfo{"127.0.0.1", "root", "huyi65", "hygame", 3306})
	fmt.Println("key is ", key)
	/*
		//DB test
		sqlpart.SqlNotQuery(key, fmt.Sprintf("insert into test values(%d,'%s',%d)", 2, "test2", 3))
		res, ok := sqlpart.SqlSelect(key, "select * from test")
		fmt.Println(res)
	*/
	/*
		// threadpool test
		t := test{
			Id:   4,
			Name: "huyi",
		}
		var arglist interface{} = t
		job1 := commondef.StJobInfo{
			RepeatTimes: 10,
			Job:         Job1,
			ArgList:     arglist,
		}
		job2 := commondef.StJobInfo{
			RepeatTimes: 20,
			Job:         Job2,
		}
		job3 := commondef.StJobInfo{
			RepeatTimes: 5,
			Job:         Job3,
		}
		threadpool.AddTask(job1)
		threadpool.AddTask(job2)
		threadpool.AddTask(job3)
	*/
	serverpart.SetCallDisConn(DisConn)
	serverpart.SetCallRead(ReadData)
	serverpart.SetIpPort("127.0.0.1", "8099")
	serverpart.StartServer()
	//	time.Sleep(10 * time.Second)
	//threadpool.StopThreadPool()
	return true
}
