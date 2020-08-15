# Session Initiation Protocol  
## Install
```sh
go get github.com/ddddddO/sip
```

## Usage  
```go

import (
    "github.com/ddddddO/sip"
)

func XXXXX() {
	availableSessions := sip.GetAvailableSessions(sip.NewConfig(
		true, laddr, clientCnt, // for SIP Server setup config.
		true, raddrs, // for SIP Client setup config.
	))

	for i := range availableSessions {
		func(ss *sip.Session) {
			// send to server
			if err := ss.Write([]byte("sending to server")); err != nil {
				panic(err)
			}

			// recieve from server
			res, err := ss.Read()
			if err != nil {
				panic(err)
			}
			log.Print(string(res))
		}(availableSessions[i])
	}
}
```

## Examples  
### Pattern in which the roles of server and client are separated.
github.com/ddddddO/sip/example/pattern1

### Pattern in which the roles of server and client are integrated.
github.com/ddddddO/sip/example/pattern2

## Reference  
- [VoIPとSIP](https://www.nic.ad.jp/ja/newsletter/No29/100.html)
- [rfc3261](https://tools.ietf.org/html/rfc3261)
- [SIPリクエスト・レスポンス例](https://tools.ietf.org/html/rfc3261#section-24.2)
- [SIPのwiki](https://ja.wikipedia.org/wiki/Session_Initiation_Protocol)
- [UDPソケットをGoで叩く](https://ascii.jp/elem/000/001/411/1411547/)
- [UDPを叩くコマンド(nc)](https://www.ecoop.net/memo/archives/udp-connection-on-command-line.html)