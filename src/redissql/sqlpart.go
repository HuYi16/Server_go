package redissql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"string"
)

//db handle 
type stDbInfo struct{
	var sDbType      string
	var SDbUserName  string
	var sDbPassWord  string
	var sDbIpPort    string
	var sDataBase    string
	var dbHandle     sql.DB
}
var dbhandle = &sql.DB{}
var bInitSql bool
func init() {
	bInitSql = false
	fmt.Printf("first of all do init\n")
}
func connectSql() bool {
	if checkSqlConnect() {
		fmt.Print("sql has connected...\n")
		return true
	}
	var err error
	dbhandle, err = sql.Open("mysql", "root:huyi65@/TradeSrc?charset=utf8")
	if err != nil {
		fmt.Printf("connect err!%s\n", err)
		return false
	}
	return true
}

func closeSql() {
	if checkSqlConnect() {
		dbhandle.Close()
		fmt.Print("close sql suc!!!!\n")
		dbhandle = nil
		bInitSql = false
		return
	}
	return
}

func checkSqlConnect() bool {
	return bInitSql
}

/*
func doSql(sql string) map[int]string,string{
	if sql == nil {
		return nil,"sql is nil!!\n "
	}
	if !checkSqlConnect(){
		return nil,"sql is not connect!!\n"
	}
	rows,err := dbhandle.Query(sql)
	if err != nil{
		return nil,err
	}
	var mRes map[int]string
	mRes = make(map[int]string)
	return mRes,nil

}
*/
