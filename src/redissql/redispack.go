package redissql

import (
	"fmt"
)

func StartRedisSql() bool {
	if checkRedisAndSqlStart() {
		fmt.Printf("redis and sql is running...\n")
		return true
	} else {
		closeSql()
	}
	if connectSql() {
		fmt.Printf("start sql suc!!!\n")
	} else {
		fmt.Printf("start sql fail!!!!\n")
		return false
	}
	fmt.Printf("start redis and sql suc!!!\n")
	return true
}

func ShutDownRedisSql() {
	fmt.Printf("Shut down redis and sql suc!!!!\n")
	closeSql()
	return
}

func checkRedisAndSqlStart() bool {
	return checkSqlConnect()
}
