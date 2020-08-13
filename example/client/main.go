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

	connectedSessionCh := sip.NewConnectedSessionCh()
	go client.Run(connectedSessionCh)

	for session := range connectedSessionCh {
		// send to server
		if err := session.Write([]byte("Hey! by client!\n")); err != nil {
			panic(err)
		}

		// recieve from server
		b, err := session.Read()
		if err != nil {
			panic(err)
		}
		log.Print(string(b))
	}
}
