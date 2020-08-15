package sip

import (
	"time"
)

type sipConfig struct {
	// for Server setup
	enableServer bool
	serverConf   serverConfig

	// for Client setup
	enableClient bool
	clientConf   clientConfig
}

type serverConfig struct {
	localAddr        string
	connectClientCnt int // 同時接続するクライアント数
}

type clientConfig struct {
	remoteAddrs []string
}

// わかりづらい気がする
func NewConfig(es bool, laddr string, cliCnt int, ec bool, raddrs []string) sipConfig {
	return sipConfig{
		enableServer: es,
		serverConf: serverConfig{
			localAddr:        laddr,
			connectClientCnt: cliCnt,
		},

		enableClient: ec,
		clientConf: clientConfig{
			remoteAddrs: raddrs,
		},
	}
}

func GetAvailableSessions(conf sipConfig) []*Session {
	// Server setup
	var sessionChByServer chan *Session
	if conf.enableServer {
		laddrUAS := conf.serverConf.localAddr
		serverUAS := NewServer(laddrUAS)
		sessionChByServer = NewConnectedSessionCh()
		clientCnt := conf.serverConf.connectClientCnt
		go serverUAS.Run(sessionChByServer, clientCnt)
	}

	// Client setup
	var sessionChByClient chan *Session
	if conf.enableClient {
		// FIXME: waiting for remote server setup...
		time.Sleep(10 * time.Second)

		clientUAC := NewClient()
		for _, raddr := range conf.clientConf.remoteAddrs {
			raddrUAS := raddr
			sessionUAS := NewSession(raddr)
			clientUAC.AddSession(raddrUAS, sessionUAS)
		}
		sessionChByClient = NewConnectedSessionCh()
		go clientUAC.Run(sessionChByClient)
	}

	// Aggregate connected session
	var availableSessions []*Session
	if conf.enableServer && conf.enableClient {
		availableSessions = aggregateAvailableSessions(
			sessionChByClient, sessionChByServer)
	} else if conf.enableServer {
		availableSessions = aggregateAvailableSessions(
			sessionChByServer)
	} else if conf.enableClient {
		availableSessions = aggregateAvailableSessions(
			sessionChByClient)
	}
	return availableSessions
}

// TODO: クライアントが持つすべてのsesssionをクローズする処理
func Close() {}

func aggregateAvailableSessions(sessionChs ...chan *Session) []*Session {
	var sss []*Session
	for i := range sessionChs {
		for session := range sessionChs[i] {
			sss = append(sss, session)
		}
	}
	return sss
}
