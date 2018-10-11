package gservergs

import (
	"commondef"
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
	go redispack.StartRedis()
	return true
}
