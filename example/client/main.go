package main

import (
	"log"

	"github.com/ddddddO/sip"
)

func main() {
	log.Print("start sip uac(client proccess)")

	raddrs := []string{"localhost:5060"}
	availableSessions := sip.GetAvailableSessions(sip.NewConfig(
		false, "", 0, // for Server setup
		true, raddrs, // for Client setup
	))

	for i := range availableSessions {
		func(ss *sip.Session) {
			// send to server
			if err := ss.Write([]byte("Hey! by client!")); err != nil {
				panic(err)
			}

			// recieve from server
			b, err := ss.Read()
			if err != nil {
				panic(err)
			}
			log.Print(string(b))
		}(availableSessions[i])
	}
}
