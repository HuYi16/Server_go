package gservergs

import (
	"fmt"
	"net/http"
	"time"
)

func update2Web() {
	iTempLock := 0
	if ServerInfo.BtempLock {
		iTempLock = 1
	}
	url := fmt.Sprintf("http://www.hycxx.top/updateserver?ip=%s&port=%s&online=%d&lock=%d&balance=%d", ServerInfo.Ip, ServerInfo.Port, ServerInfo.NowNumber, iTempLock, ServerInfo.BalanceNumber)
	_, err := http.Get(url)
	//	fmt.Println(url)
	if err != nil {
		//		fmt.Println("update err!!!", err.Error())
	} else {
		//		fmt.Println("update suc!!!", res)
	}
}

func StartTimer() {
	updateTimer := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-updateTimer.C:
			update2Web()
			if ServerInfo.BstopServerFlag {
				break
			}
		}
	}
}
