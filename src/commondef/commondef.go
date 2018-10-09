// all defines struct enum

package commondef

//server
type StServerInfo struct {
	Ip             string //ip of server
	Port           string //port info
	NowNumber      int32  //online user of server
	BalanceNumber  int32  //the number of need balance for GS
	LastUpdateTime int64  //the time of last update
	BtempLock      bool   //the server is on
}
