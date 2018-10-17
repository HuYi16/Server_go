package sqlpart

import (
	"commondef"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var mapDBIndex map[string]*sql.DB

func init() {
	fmt.Println("init sql part")
	mapDBIndex = make(map[string]*sql.DB)
}

func connctDB(key string, arg commondef.StSqlRedisBaseInfo) bool {
	sqlName := "mysql"
	/*
		if len(mapDBIndex) >= 1 {
			sqlName = fmt.Sprintf("mysql%d", len(mapDBIndex))
			sql.Register(sqlName, &sql.MySQLDriver{})
		}
	*/
	info := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", arg.DBUser, arg.DBPassWord, arg.DBHost, arg.DBPort, arg.DBName)
	db, err := sql.Open(sqlName, info)
	fmt.Println("err:", err, "info:", info)
	if err == nil {
		fmt.Println("connect db ", info, " suc!!!")
		mapDBIndex[key] = db
		return true
	}
	fmt.Println("connect db ", info, " fail", err)
	return false
}
func checkBaseInfo() bool {
	if len(mapDBIndex) == 0 {
		return false
		//	mapHostInfo["127.0.0.1roothuyi65hygame"] = &commondef.StSqlRedisBaseInfo{"127.0.0.1", "root", "huyi65","hygame"}
	}
	/*
		if len(mapDBIndex) != len(mapHostInfo) {
			for k, v := range mapHostInfo {
				_, ok := mapDBIndex[k]
				if !ok {
					return connectDB(key,v.DBHost, v.DBUser, v.DBPassWord,v.DBName)
				}
			}
		}
	*/
	return true
}

func StartSql(arg commondef.StSqlRedisBaseInfo) string {
	if arg.DBHost == "" || arg.DBUser == "" || arg.DBName == "" || arg.DBPort <= 0 {
		return ""
	}
	key := fmt.Sprintf("%s%s%s%s%d", arg.DBHost, arg.DBUser, arg.DBPassWord, arg.DBName, arg.DBPort)
	_, ok := mapDBIndex[key]
	if !ok {
		fmt.Println(key)
		if !connctDB(key, arg) {
			return ""
		}
	}
	return key
}

func CloseDB(key string) {
	v, ok := mapDBIndex[key]
	if ok {
		v.Close()
		delete(mapDBIndex, key)
	}
}
func EndSQL() {
	for k, v := range mapDBIndex {
		v.Close()
		delete(mapDBIndex, k)
	}
}

func SqlNotQuery(index, sql string) bool {
	if sql == "" {
		return false
	}
	v, ok := mapDBIndex[index]
	if !ok {
		fmt.Println("key is not exit!!", index)
		return false
	}
	handle, err := v.Prepare(sql)
	if err != nil {
		panic(err)
		return false
	}
	_, errexec := handle.Exec()
	if errexec != nil {
		panic(errexec)
		return false
	}
	return true
}

func SqlSelect(index, sql string) (map[int]map[string]string, bool) {
	res := make(map[int]map[string]string)
	if sql == "" {
		return res, false
	}
	v, ok := mapDBIndex[index]
	if !ok {
		fmt.Println("key is not exit!!", index)
		return res, false
	}
	rows, _ := v.Query(sql)
	cols, _ := rows.Columns()
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	for rows.Next() {
		rows.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		res[i] = row
		i++
	}
	return res, true
}
