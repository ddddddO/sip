package sip

import (
	"bufio"
	"errors"
	"net"
)

func NewConnectedSessionCh() chan *Session {
	return make(chan *Session)
}

// Key: remote address / Value: Session
type SessionMap map[string]*Session

type Session struct {
	raddr *net.UDPAddr
	conn  *net.UDPConn
	br    *bufio.Reader
	bw    *bufio.Writer

	state      State
	originType OriginType
}

// sessionをクライアント・サーバーのどちらで生成したかの判断材料
// ss.Write/ss.Read内で利用するため
type OriginType string

const (
	OriginServer = OriginType("server")
	OriginClient = OriginType("client")
)

// 一旦、クライアントから接続用
func NewSession(raddr string) *Session {
	conn, err := net.Dial("udp4", raddr)
	if err != nil {
		panic(err)
	}
	udpConn, ok := conn.(*net.UDPConn)
	if !ok {
		panic("can't assertion")
	}
	resolvedRaddr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		panic(err)
	}

	return &Session{
		raddr: resolvedRaddr,
		conn:  udpConn,
		br:    bufio.NewReader(udpConn),
		bw:    bufio.NewWriter(udpConn),

		state:      StateINIT,
		originType: OriginClient,
	}
}

func (ss *Session) GetState() State {
	return ss.state
}

func (ss *Session) ChangeState(state State) {
	ss.state = state
}

func (ss *Session) Read(buf []byte) (int, error) {
	if ss.originType == OriginClient {
		if ss.br == nil {
			return 0, errors.New("can't access")
		}

		n, err := ss.br.Read(buf)
		if err != nil {
			return n, err
		}
		return n, nil
	} else if ss.originType == OriginServer {
		n, _, err := ss.conn.ReadFrom(buf)
		if err != nil {
			return n, err
		}
		return n, nil
	}
	return 0, errors.New("undefined origin type")
}

func (ss *Session) Write(b []byte) (int, error) {
	if ss.originType == OriginClient {
		if ss.bw == nil {
			return 0, errors.New("can't access")
		}
		n, err := ss.bw.Write(b)
		if err != nil {
			return n, err
		}
		if err := ss.bw.Flush(); err != nil {
			return n, err
		}
		return n, nil
	} else if ss.originType == OriginServer {
		n, err := ss.conn.WriteTo(b, ss.raddr)
		return n, err
	}
	return 0, errors.New("undefined origin type")
}

func (ss *Session) Close() error {
	return ss.conn.Close()
}
