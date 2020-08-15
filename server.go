package sip

import (
	"net"
)

type Server struct {
	Conn  *net.UDPConn
	ssmap SessionMap
}

func NewServer(laddr string) *Server {
	resolvedLaddr, err := net.ResolveUDPAddr("udp", laddr)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", resolvedLaddr)
	if err != nil {
		panic(err)
	}

	return &Server{
		Conn:  conn,
		ssmap: SessionMap{},
	}
}

func (s *Server) existSession(raddr string) bool {
	_, ok := s.ssmap[raddr]
	return ok
}

func (s *Server) AddSession(ra *net.UDPAddr) {
	raddr := ra.String()
	if !s.existSession(raddr) {
		// NOTE: ここのあたりが微妙
		s.ssmap[raddr] = &Session{
			raddr: ra,
			conn:  s.Conn,
			br:    nil,
			bw:    nil,

			state:      StateINIT,
			originType: OriginServer,
		}
	}
}

// FIXME: clientCnt: 苦肉の策。綺麗ではないし、不便
func (s *Server) Run(connectedSessionCh chan<- *Session, clientCnt int) error {
	var connectedCnt int
	for {
		if connectedCnt == clientCnt {
			close(connectedSessionCh)
			break
		}

		b := make([]byte, 1024)
		_, ra, err := s.Conn.ReadFromUDP(b)
		if err != nil {
			return err
		}
		raddr := ra.String()
		s.AddSession(ra)

		LogInfo(raddr, s.ssmap[raddr])

		switch s.ssmap[raddr].GetState() {
		case StateINIT:
			s.ssmap[raddr].ChangeState(StateRINGING)

			ringingRes := buildResponseRINGING()
			if _, err := s.ssmap[raddr].Write(ringingRes); err != nil {
				return err
			}
			//err = s.ssmap[raddr].Write([]byte("Hello?")) // NOTE: panic: write udp 127.0.0.1:5060: write: destination address required

			if isValidRequestINVITE(b) {
				okRes := buildResponseOK()
				s.ssmap[raddr].Write(okRes)
				s.ssmap[raddr].ChangeState(StateOK)
			} else {
				s.ssmap[raddr].Write([]byte("response code 4XX\r\n")) // TODO: 4XX response
			}
		case StateRINGING:
			if isValidRequestINVITE(b) {
				okRes := buildResponseOK()
				s.ssmap[raddr].Write(okRes)
				s.ssmap[raddr].ChangeState(StateOK)
			} else {
				s.ssmap[raddr].Write([]byte("response code 4XX\r\n")) // TODO: 4XX response
			}
		case StateOK:
			if isValidRequestACK(b) {
				s.ssmap[raddr].ChangeState(StateCONNECTED)
				LogInfo(raddr, s.ssmap[raddr])
				connectedCnt++
				connectedSessionCh <- s.ssmap[raddr]
			} else {
				s.ssmap[raddr].Write([]byte("response code 4XX\r\n")) // TODO: 4XX response
			}
		case StateCONNECTED:
			// TODO: ... when BYE req ?
			LogInfo(raddr, s.ssmap[raddr])
		}
	}
	return nil
}
