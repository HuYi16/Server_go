package threadpool

import (
	"commondef"
	"fmt"
	"runtime"
	"time"
)

var stPoolInfo *commondef.StPoolData

func checkPoolInfo() {
	if stPoolInfo == nil {
		stPoolInfo = &commondef.StPoolData{
			MaxNum:   runtime.NumCPU()*2 + 1,
			StopFlag: false,
		}
		//stPoolInfo.MaxNum = runtime.NumCPU()*2 + 1
		stPoolInfo.MaxJobQueue = stPoolInfo.MaxNum * 5
		//stPoolInfo.StopFlag = false
	}
}

func SetThreadNum(iMaxNum int) bool {
	checkPoolInfo()
	if iMaxNum <= 0 {
		return false
	}
	stPoolInfo.MaxNum = iMaxNum
	stPoolInfo.MaxJobQueue = iMaxNum * 5
	return true
}

func threadWork(id int) {
	checkPoolInfo()
	llStarttime := time.Now().UnixNano()
	fmt.Println("threadwork will start for ", id)
	for !stPoolInfo.StopFlag {
		llStarttime = time.Now().UnixNano()
		task, ok := <-stPoolInfo.JobQueue
		if ok {
			if task.RepeatTimes > 0 || task.RepeatTimes == -1 {
				task.Job(task.ArgList)
				if task.RepeatTimes == -1 || task.RepeatTimes-1 > 0 {
					if task.RepeatTimes > 0 {
						task.RepeatTimes--
					}
					stPoolInfo.JobQueue <- task
				}
			}
		}
		iBlank := time.Now().UnixNano() - llStarttime
		if iBlank < 2*1000000 {
			time.Sleep(time.Duration(2*1000000 - iBlank))
		}
	}
	fmt.Println("thread quit!!:", id)
}

func StartThreadPool() bool {
	checkPoolInfo()
	if stPoolInfo.MaxNum <= 0 {
		return false
	}
	stPoolInfo.StopFlag = false
	stPoolInfo.JobQueue = make(chan commondef.StJobInfo)
	for i := 0; i < stPoolInfo.MaxNum; i++ {
		go threadWork(i)
		fmt.Println("start threadwork:", i, stPoolInfo.MaxNum)
	}
	return true
}

func StopThreadPool() {
	checkPoolInfo()
	stPoolInfo.StopFlag = true
}
func AddTask(task commondef.StJobInfo) bool {
	checkPoolInfo()
	if task.RepeatTimes == 0 || task.Job == nil {
		return false
	}
	stPoolInfo.JobQueue <- task
	return true
}
