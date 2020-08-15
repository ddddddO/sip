package main

import (
	"io"
	"log"
	"os"

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
			if _, err := ss.Write([]byte("Hello! by server..\n")); err != nil {
				panic(err)
			}

			// recieve from client
			io.Copy(os.Stdout, ss)
		}(availableSessions[i])
	}
}
