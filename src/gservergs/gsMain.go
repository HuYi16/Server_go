package gservergs

import (
	"commondef"
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
	return true
}
