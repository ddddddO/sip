package main

import (
	"log"
	"sync"

	"github.com/ddddddO/sip"
)

func main() {
	laddr := "localhost:5062"
	clientCnt := 1
	raddrs := []string{"localhost:5061"}

	availableSessions := sip.GetAvailableSessions(sip.NewConfig(
		true, laddr, clientCnt, // for Server setup
		true, raddrs, // for Client setup
	))

	// send & recieve msg!
	wg := &sync.WaitGroup{}
	for i := range availableSessions {
		wg.Add(1)
		go func(ss *sip.Session) {
			if _, err := ss.Write([]byte("Bye by ua_2")); err != nil {
				panic(err)
			}

			res := make([]byte, 1024)
			_, err := ss.Read(res)
			if err != nil {
				panic(err)
			}
			log.Print("recieve msg:", string(res))
			wg.Done()
		}(availableSessions[i])
	}
	wg.Wait()
}
