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
			state: State("init"),
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
		s.AddSession(raddr)

		log.Print(string(b))

		if isValidINVITERequest(b) {
			s.Conn.WriteTo([]byte("response code 2XX"), ra)
		} else {
			s.Conn.WriteTo([]byte("response code 4XX"), ra)
		}

		//err = s.ssmap[raddr].Write([]byte("Hello?")) // NOTE: panic: write udp 127.0.0.1:5060: write: destination address required
		_, err = s.Conn.WriteTo([]byte("Hello?"), ra)
		if err != nil {
			panic(err)
		}
	}
}

var (
	end    = []byte("\r\n")
	method = []byte("INVITE")
)

// NOTE: 一旦ヘッダーが妥当かチェック
func isValidINVITERequest(b []byte) bool {
	requestLine := bytes.Split(b, end)[0]
	splited := bytes.Split(requestLine, []byte(" "))

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
