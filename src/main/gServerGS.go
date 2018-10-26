package main

import (
	"fmt"
	"gservergs"
	//"time"
)

func init() {
	fmt.Println("gs main init!!!")
}

func main() {
	if gservergs.StartGs() {
		fmt.Println("StartGs suc!!!")
	} else {
		fmt.Println("StartGs fail!!!")
	}
	fmt.Println("end !!!!!!")
	//	time.Sleep(10000 * time.Second)
}
