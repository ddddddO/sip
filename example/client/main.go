package main

import (
	"log"

	"github.com/ddddddO/sip"
)

func main() {
	log.Print("start sip uac(client proccess)")
	raddr := "localhost:5060"
	session := sip.NewSession(raddr)
	defer session.Close()

	client := sip.NewClient()
	client.AddSession(raddr, session)

	err := client.Run()
	if err != nil {
		panic(err)
	}

}
