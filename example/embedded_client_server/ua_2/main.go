package main

import (
	"log"
	"sync"
	"time"

	"github.com/ddddddO/sip"
)

func main() {
	// Server setup
	laddrUAS1 := "localhost:5062"
	serverUAS1 := sip.NewServer(laddrUAS1)
	sessionChByServer := sip.NewConnectedSessionCh()
	clientCnt := 1
	go serverUAS1.Run(sessionChByServer, clientCnt)

	// FIXME: waiting for remote server setup...
	time.Sleep(10 * time.Second)

	// Client setup
	raddrUAS2 := "localhost:5061"
	sessionUAS2 := sip.NewSession(raddrUAS2)
	defer sessionUAS2.Close()
	clientUAC1 := sip.NewClient()
	clientUAC1.AddSession(raddrUAS2, sessionUAS2)
	sessionChByClient := sip.NewConnectedSessionCh()
	go clientUAC1.Run(sessionChByClient)

	// Aggregate connected session
	availableSessions := []*sip.Session{}
	for session := range sessionChByServer {
		availableSessions = append(availableSessions, session)
	}
	for session := range sessionChByClient {
		availableSessions = append(availableSessions, session)
	}

	// send & recieve msg!
	wg := &sync.WaitGroup{}
	for i := range availableSessions {
		wg.Add(1)
		go func(ss *sip.Session) {
			ss.Write([]byte("HeyHey! by ua_2"))

			res, err := ss.Read()
			if err != nil {
				panic(err)
			}
			log.Print("recieve msg:", string(res))
			wg.Done()
		}(availableSessions[i])
	}
	wg.Wait()
}
