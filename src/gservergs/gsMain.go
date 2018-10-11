package gservergs

import (
	"commondef"
	//	"fmt"
	"redispack"
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

func StartGs() bool {
	go StartTimer()
	go redispack.StartRedis2()
	//	redispack.RedisSet(1, "test", "suc")
	//	redispack.RedisSet(1, "test1", "succ")
	//	ok, val := redispack.RedisGet(1, "test")
	//	fmt.Println("test", ok, val)
	return true
}
