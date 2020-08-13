package sip

import (
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

			ringingRes := buildResponseRINGING()
			s.Conn.WriteTo(ringingRes, ra)

			if isValidRequestINVITE(b) {
				okRes := buildResponseOK()
				s.Conn.WriteTo(okRes, ra)
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
