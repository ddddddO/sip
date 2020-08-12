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

		//err = s.ssmap[raddr].Write([]byte("Hello?")) // NOTE: panic: write udp 127.0.0.1:5060: write: destination address required
		_, err = s.Conn.WriteTo([]byte("Hello?"), ra)
		if err != nil {
			panic(err)
		}
	}
}
