// all defines struct enum

package commondef

//server
type StServerInfo struct {
	Ip              string //ip of server
	Port            string //port info
	NowNumber       int32  //online user of server
	BalanceNumber   int32  //the number of need balance for GS
	LastUpdateTime  int64  //the time of last update
	BtempLock       bool   //the server is on
	BstopServerFlag bool   //the flag of server will stop
}

//threadpool

type StJobInfo struct {
	RepeatTimes int
	Job         func()
}

type StPoolData struct {
	MaxNum      int
	StopFlag    bool
	MaxJobQueue int
	JobQueue    chan StJobInfo
}
