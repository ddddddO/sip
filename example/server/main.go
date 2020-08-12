package main

import (
	"log"

	"github.com/ddddddO/sip"
)

func main() {
	log.Print("start sip uas(server proccess)")

	laddr := "localhost:5060"
	server := sip.NewServer(laddr)

	for {
		raddr, buf, err := server.Read()
		if err != nil {
			panic(err)
		}
		log.Print(string(buf))

		if err := server.Write([]byte("Hello!!"), raddr); err != nil {
			panic(err)
		}
	}
}
