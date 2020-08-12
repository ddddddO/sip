package main

import (
	"log"

	"github.com/ddddddO/sip"
)

func main() {
	log.Print("start sip uas(server proccess)")

	laddr := "localhost:5060"
	server := sip.NewServer(laddr)

	if err := server.Run(); err != nil {
		panic(err)
	}
}
