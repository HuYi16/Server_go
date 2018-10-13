package sqlpack

import (
	"commondef"
	"fmt"
)

var hostInfo *commondef.StSqlRedisBaseInfo

func init() {
	fmt.Println("init sql part")
}

func checkBaseInfo() {
	if nil == hostInfo {
		hostInfo = &commondef.StSqlRedisBaseInfo{"127.0.0.1", "root", "huyi65"}
	}
}

func SetSqlBaseInfo(host, user, psw string) bool {
	if host == "" || user == "" {
		return false
	}
	if nil == hostInfo {
		fmt.Println(host, user, psw)
		hostInfo = &commondef.StSqlRedisBaseInfo{host, user, psw}
	}
	return true
}
