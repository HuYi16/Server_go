package gserverweb

import (
	"bufio"
	"bytes"
	"commondef"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

/*
type gsInfo struct {
	IpAndPort      string //ip of server
	Port           string //port of server
	NowNumber      int32  // online user of server
	LastUpdateTime int64  // the time of last uodate this server
}
*/

var gsInfoMap map[string]commondef.StServerInfo

func init() {
	gsInfoMap = make(map[string]commondef.StServerInfo)
	if nil == gsInfoMap {
		fmt.Println("map is nil")
	}
}

// func updateserver online info
func Updateserver(w http.ResponseWriter, r *http.Request) {
	//	key,iOnlineNumber := getIpAndPort("")
	r.ParseForm()
	sIp := r.FormValue("ip")
	sNumber := r.FormValue("online")
	sPort := r.FormValue("port")
	sTempLock := r.FormValue("lock")
	bTempLockflag := false
	sBalanceNumber := r.FormValue("balance")
	if sIp == "" || sNumber == "" || sPort == "" || sTempLock == "" {
		fmt.Fprintf(w, "update param key is empty!!")
		return
	}
	iOnlineNumber, okatoi := strconv.Atoi(sNumber)
	if okatoi != nil {
		iOnlineNumber = 0
	}
	iBalanceNumber, okbalance := strconv.Atoi(sBalanceNumber)
	if okbalance != nil {
		iBalanceNumber = 100
	}
	iTempLockFlag, okflag := strconv.Atoi(sTempLock)
	if okflag != nil {
		iTempLockFlag = 0
	}
	if iTempLockFlag != 0 {
		bTempLockflag = true
	}
	gsInfoMap[sIp+":"+sPort] = commondef.StServerInfo{Ip: sIp, Port: sPort, LastUpdateTime: time.Now().Unix(), NowNumber: int32(iOnlineNumber), BtempLock: bTempLockflag, BalanceNumber: int32(iBalanceNumber)}
	fmt.Fprintf(w, "update suc!!")
	return
}

// clent gey server IpAndPort and delete offline server which updateinfo befor 15s
func Getserver(w http.ResponseWriter, r *http.Request) {
	if len(gsInfoMap) == 0 {
		fmt.Fprintf(w, "gs list is empty!!!pls wait...")
		return
	}
	sTempKey := ""
	iUpNumber := 0
	for k, v := range gsInfoMap {
		t := time.Now()
		if t.Unix()-v.LastUpdateTime > 15 {
			delete(gsInfoMap, k)
		}
		if !v.BtempLock {
			if v.BalanceNumber >= v.NowNumber {
				fmt.Fprintf(w, "%s", k)
				return
			} else {
				if sTempKey == "" {
					sTempKey = k
					iUpNumber = int(v.NowNumber - v.BalanceNumber)
				} else if iUpNumber < int(v.NowNumber-v.BalanceNumber) {
					sTempKey = k
					iUpNumber = int(v.NowNumber - v.BalanceNumber)
				}
			}
		}
	}
	fmt.Fprintf(w, "%s", sTempKey)
	return
}

func GetAllServer(w http.ResponseWriter, r *http.Request) {
	sRes := ""
	for _, v := range gsInfoMap {
		buf := bytes.NewBuffer(make([]byte, 80))
		bw := bufio.NewWriter(buf)
		fmt.Fprintf(bw, "ip:%s port:%s online:%d balance:%d lockflag:%t\n", v.Ip, v.Port, v.NowNumber, v.BalanceNumber, v.BtempLock)
		bw.Flush()
		sRes += buf.String()
	}
	fmt.Fprintf(w, "serverlist:\n%s", sRes)
	return
}
