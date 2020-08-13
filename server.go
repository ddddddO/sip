package sip

import (
	"bytes"
	"log"
	"net"
)

type Server struct {
	Conn  *net.UDPConn
	ssmap SessionMap
}

func NewServer(laddr string) *Server {
	resolvedLaddr, err := net.ResolveUDPAddr("udp", laddr)
	conn, err := net.ListenUDP("udp", resolvedLaddr)
	if err != nil {
		panic(err)
	}

	return &Server{
		Conn:  conn,
		ssmap: SessionMap{},
	}
}

func (s *Server) AddSession(raddr string) {
	if _, ok := s.ssmap[raddr]; !ok {
		// NOTE: ここのあたりが微妙
		s.ssmap[raddr] = &Session{
			br:    nil,
			bw:    nil,
			state: StateINIT,
		}
	}
}

func (s *Server) Run() error {
	for {
		b := make([]byte, 1024)
		_, ra, err := s.Conn.ReadFrom(b)
		if err != nil {
			return err
		}
		raddr := ra.String()
		s.AddSession(raddr) // TODO: 既にセッションが存在するかチェックを関数に切り出すか考える

		log.Printf("debug\n%s", string(b))

		switch s.ssmap[raddr].GetState() {
		case StateINIT:
			s.ssmap[raddr].ChangeState(StateRINGING)

			// TODO: Add 180 Ringing response to uac

			if isValidINVITERequest(b) {
				s.Conn.WriteTo([]byte("response code 2XX"), ra)
				s.ssmap[raddr].ChangeState(StateOK)
			} else {
				s.Conn.WriteTo([]byte("response code 4XX"), ra)
			}
		case StateOK:
			// if ack?
			// true -> ChangeState(connected)
		}

		//err = s.ssmap[raddr].Write([]byte("Hello?")) // NOTE: panic: write udp 127.0.0.1:5060: write: destination address required
		_, err = s.Conn.WriteTo([]byte("Hello?"), ra)
		if err != nil {
			panic(err)
		}
	}
}

var (
	space  = []byte(" ")
	end    = []byte("\r\n")
	method = []byte("INVITE")
)

func isValidINVITERequest(b []byte) bool {
	requestMsg := bytes.Split(b, end)
	// check request line
	if !isValidRequestLine(requestMsg[0]) {
		return false
	}

	// check request header
	if !isValidRequestHeader(requestMsg[1:]) {
		return false
	}

	return true
}

func isValidRequestLine(requestLine []byte) bool {
	splited := bytes.Split(requestLine, space)
	if len(splited) != 3 {
		return false
	}
	// check Method
	if !bytes.Equal(splited[0], method) {
		return false
	}
	// TODO:check Request-URI

	// TODO:check SIP-Version

	return true
}

func isValidRequestHeader(requestHeaders [][]byte) bool {
	// TODO: check Via/To/From/Call-ID/CSeq/Contact/Max-Forwards/Content-Type/Content-Length
	// 「This example contains a minimum required set.」最低限のヘッダのよう
	return true
}
