package main

import (
	"log"

	"github.com/ddddddO/sip"
)

func main() {
	log.Print("start sip uas(server proccess)")

	laddr := "localhost:5060"
	server := sip.NewServer(laddr)
	connectedSessionCh := sip.NewConnectedSessionCh()

	// TODO: 複数の時の対応
	clientCnt := 2 // 接続してくるクライアントの数を指定しないといけない。。。
	go server.Run(connectedSessionCh, clientCnt)

	bufSessions := []*sip.Session{}
	for session := range connectedSessionCh {
		bufSessions = append(bufSessions, session)
	}

	for i := range bufSessions {
		func(session *sip.Session) {
			// send to client
			if err := session.Write([]byte("Hello! by server..\n")); err != nil {
				panic(err)
			}

			// recieve from client
			b, err := session.Read()
			if err != nil {
				panic(err)
			}
			log.Print(string(b))
		}(bufSessions[i])
	}
}
