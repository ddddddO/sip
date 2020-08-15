package main

import (
	"log"

	"github.com/ddddddO/sip"
)

func main() {
	log.Print("start sip uas(server proccess)")

	laddr := "localhost:5060"
	clientCnt := 1
	availableSessions := sip.GetAvailableSessions(sip.NewConfig(
		true, laddr, clientCnt, // for Server setup
		false, nil, // for Client setup
	))

	for i := range availableSessions {
		func(ss *sip.Session) {
			// send to client
			if err := ss.Write([]byte("Hello! by server..")); err != nil {
				panic(err)
			}

			// recieve from client
			b, err := ss.Read()
			if err != nil {
				panic(err)
			}
			log.Print(string(b))
		}(availableSessions[i])
	}
}
