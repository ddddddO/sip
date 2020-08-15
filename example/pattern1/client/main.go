package main

import (
	"io"
	"log"
	"os"

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
			if _, err := ss.Write([]byte("Hey! by client!\n")); err != nil {
				panic(err)
			}

			// recieve from server
			io.Copy(os.Stdout, ss)
		}(availableSessions[i])
	}
}
