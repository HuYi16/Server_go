package main

import (
	"gservergs"
	"log"
)

var stServerInfo ServerInfo

func init() {
	log.Info("gs main init")
}

func main() {
	if gservergs.StartGs() {
		log.Info("start gs suc")
	} else {
		log.Info("start gs fail")
	}
}
