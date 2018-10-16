package redispack

import (
	"commondef"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"net/http"
	"runtime"
	"time"
)

var hostInfo *commondef.StSqlRedisBaseInfo

//redis pool method one
var MAX_POOL_SIZE = 6
var redisPool chan redis.Conn

func putRedis(conn redis.Conn) {
	if redisPool == nil {
		redisPool = make(chan redis.Conn, MAX_POOL_SIZE)
	}
	if len(redisPool) >= MAX_POOL_SIZE {
		conn.Close()
		return
	}
	redisPool <- conn
}

func initRedis(network, address string) redis.Conn {
	if len(redisPool) == 0 {
		fmt.Println("initredis")
		redisPool = make(chan redis.Conn, MAX_POOL_SIZE)
		go func() {
			for i := 0; i < MAX_POOL_SIZE/2; i++ {
				c, err := redis.Dial(network, address)
				fmt.Println("make conn", i, network, address)
				if err != nil {
					panic(err)
				}
				defer c.Close()
				putRedis(c)
			}
		}()
	}
	return <-redisPool
}

func redisServer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sType := r.FormValue("type")
	sKey := r.FormValue("key")
	sValue := r.FormValue("value")
	if sKey == "" || sType == "" {
		io.WriteString(w, "key  or type is empty!")
		return
	}
	//c := initRedis("tcp", "47.106.141.213:6379")
	c := initRedis("tcp", "127.0.0.1:6379")
	if sType == "LPUSH" || sType == "SET" {
		if sValue == "" {
			io.WriteString(w, fmt.Sprintf("value is empty!!key:%s,type:%s", sKey, sType))
			return
		}
		fmt.Println("redispool len:", len(redisPool))
		_, err := c.Do(sType, sKey, sValue)
		if err != nil {
			io.WriteString(w, "suc!!")
			return
		} else {
			io.WriteString(w, fmt.Sprintf("fail:%s", err.Error()))
			return
		}
	} else if sType == "GET" {
		sValue, err := c.Do(sType, sKey)
		if err != nil {
			io.WriteString(w, fmt.Sprintf("%s", sValue))
			return
		} else {
			io.WriteString(w, "")
			return
		}
	}
	io.WriteString(w, "")
	return
}

func StartRedis() bool {
	fmt.Println("init redispool")
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/", redisServer)
	http.ListenAndServe(":8099", nil)
	return true

	//	return initRedisPool2()
}

//redispool second method

var (
	redisPool2 *redis.Pool
)

func initRedisPool2() {
	fmt.Println("start init pool2")
	if nil == hostInfo {
		hostInfo = &commondef.StSqlRedisBaseInfo{
			DBHost: "127.0.0.1:6379",
		}
	}
	redisPool2 = &redis.Pool{
		MaxIdle:     1,
		MaxActive:   10,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", hostInfo.DBHost)
			if err != nil {
				panic(err)
				return nil, err

			}
			//select db
			//c.Do("SELECT",db)
			return c, nil
		},
	}
}

func checkPool2() {
	if nil == redisPool2 {
		initRedisPool2()
	}
}
func RedisSet(dbid int, key string, value string) bool {
	if dbid <= 0 || key == "" || value == "" {
		return false
	}
	checkPool2()
	c := redisPool2.Get()
	defer c.Close()
	//	c.Do("SELECT",dbid)
	dbid = 0
	_, err := c.Do("SET", key, value)
	return nil == err
}

func RedisGet(dbid int, key string) (bool, string) {
	if dbid <= 0 || key == "" {
		return false, "para err"
	}
	checkPool2()
	c := redisPool2.Get()
	defer c.Close()
	dbid = 0
	value, err := redis.String(c.Do("GET", key))
	//	fmt.Println(value, err)
	return err == nil, value
}

func RedisDelete(dbid int, key string) bool {
	if dbid <= 0 || key == "" {
		return false
	}
	checkPool2()
	c := redisPool2.Get()
	defer c.Close()
	dbid = 0
	_, err := c.Do("DELETE", key)
	return err == nil
}

func SetRedisBaseInfo(host, user, psw string) bool {
	if host == "" {
		return false
	}
	fmt.Println("redis set host", host, user, psw)
	if nil == hostInfo {
		hostInfo = &commondef.StSqlRedisBaseInfo{host, user, psw, "", 0}
	}
	return true
}
