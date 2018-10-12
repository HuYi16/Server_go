package gservergs

import (
	"commondef"
	//	"fmt"
	//	"redispack"
	"threadpool"
	//	"time"
)

var ServerInfo commondef.StServerInfo

func LoadConfig() {
	ServerInfo.Ip = "47.106.141.213"
	ServerInfo.Port = "8001"
	ServerInfo.NowNumber = 0
	ServerInfo.BalanceNumber = 100
	ServerInfo.BtempLock = false
}
func init() {
	LoadConfig()
}

/*
func Job1() {
	fmt.Println("job1 test")
	time.Sleep(1 * time.Second)
}
func Job2() {
	fmt.Println("job2 test")
	time.Sleep(3 * time.Second)
}
func Job3() {
	fmt.Println("job3 test")
	time.Sleep(2 * time.Second)
}*/
func StartGs() bool {
	go StartTimer()
	//	redispack.RedisSet(1, "test", "suc")
	//	redispack.RedisSet(1, "test1", "succ")
	//	ok, val := redispack.RedisGet(1, "test")
	//	fmt.Println("test", ok, val)
	//	threadpool.SetThreadNum(3)
	threadpool.StartThreadPool()
	/*	job1 := commondef.StJobInfo{
			RepeatTimes: 10,
			Job:         Job1,
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
	//	time.Sleep(10 * time.Second)
	//threadpool.StopThreadPool()
	return true
}
