package gservergs

import (
	"commondef"
	"log"
)

var ServerInfo commondef.ServerInfo

func init() {
	log.Info("gsservergs init!")
	ServerInfo.OnlineNumber = 0
}

func StartGs() bool {
	log.Info("start GS suc!!!")
	return true
}
