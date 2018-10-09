package main

import (
	"fmt"
	"gserverweb"
	"log"
	"net/http"
)

//register all func
func register() {
	http.HandleFunc("/getserver", gserverweb.Getserver)
	http.HandleFunc("/updateserver", gserverweb.Updateserver)
	http.HandleFunc("/getallonlineserver", gserverweb.GetAllServer)
}

// start web service
func startService() int {
	register()
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe err:", err)
		return 1
	}
	return 0
}
func main() {
	if 0 != startService() {
		fmt.Println("start webService err!!")
		return
	}
}
