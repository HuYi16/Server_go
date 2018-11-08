package main

import (
	"fmt"
	"gserverls"
)

func init() {
	fmt.Println("login server init!!!")
}

func main() {
	fmt.Println("login server start!!!")
	if gserverls.StartLS() {
		fmt.Println("start ls suc!!")
	} else {
		fmt.Println("start ls fail!!")
	}
	return
}
