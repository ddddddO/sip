# Session Initiation Protocol  
## Install
```sh
go get github.com/ddddddO/sip
```

## Usage  
### Server  
```go

import (
	"io"
	"log"
	"os"

	"github.com/ddddddO/sip"
)

func XXXXX() {
	enableServer := true
	enableClient := false
	laddr := "localhost:5060"
	clientCnt := 1
	availableSessions := sip.GetAvailableSessions(sip.NewConfig(
		enableServer, laddr, clientCnt, // for Server setup
		enableClient, nil, // for Client setup
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
```

### Client  

```go
import (
	"io"
	"log"
	"os"

	"github.com/ddddddO/sip"
)

func YYYYY() {
	enableServer := false
	enableClient := true
	raddrs := []string{"localhost:5060"}
	availableSessions := sip.GetAvailableSessions(sip.NewConfig(
		enableServer, "", 0, // for Server setup
		enableClient, raddrs, // for Client setup
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
```

## Examples  
### Pattern in which the roles of server and client are separated.
github.com/ddddddO/sip/example/pattern1

### Pattern in which the roles of server and client are integrated.
github.com/ddddddO/sip/example/pattern2

## Reference  
- [SIPについて詳しい・わかりやすい](https://www.atmarkit.co.jp/ait/articles/0711/16/news146.html)  
- [VoIPとSIP](https://www.nic.ad.jp/ja/newsletter/No29/100.html)
- [rfc3261](https://tools.ietf.org/html/rfc3261)
- [SIPリクエスト・レスポンス例](https://tools.ietf.org/html/rfc3261#section-24.2)
- [SIPのwiki](https://ja.wikipedia.org/wiki/Session_Initiation_Protocol)
- [UDPソケットをGoで叩く](https://ascii.jp/elem/000/001/411/1411547/)
- [UDPを叩くコマンド(nc)](https://www.ecoop.net/memo/archives/udp-connection-on-command-line.html)