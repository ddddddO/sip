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

			ringingRes := s.buildResponseRINGING()
			s.Conn.WriteTo(ringingRes, ra)

			if isValidRequestINVITE(b) {
				okRes := s.buildResponseOK()
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

func (s *Server) buildResponseRINGING() []byte {
	var b []byte
	b = append(b, []byte("SIP/2.0 180 Ringing\r\n")...) // 1xx: Provisional -- request received, continuing to process the request;
	b = append(b, []byte("Via: SIP/2.0/UDP server10.biloxi.com;branch=z9hG4bK4b43c2ff8.1;received=192.0.2.3\r\n")...)
	b = append(b, []byte("To: Bob <sip:bob@biloxi.com>;tag=a6c85cf\r\n")...)
	b = append(b, []byte("From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n")...)
	b = append(b, []byte("Call-ID: a84b4c76e66710\r\n")...)
	b = append(b, []byte("Contact: <sip:bob@192.0.2.4>\r\n")...)
	b = append(b, []byte("CSeq: 314159 INVITE\r\n")...)
	b = append(b, []byte("Content-Length: 0\r\n")...)

	return b
}

func (s *Server) buildResponseOK() []byte {
	var b []byte
	b = append(b, []byte("SIP/2.0 200 OK\r\n")...)
	b = append(b, []byte("Via: SIP/2.0/UDP server10.biloxi.com;branch=z9hG4bK4b43c2ff8.1;received=192.0.2.3\r\n")...)
	b = append(b, []byte("To: Bob <sip:bob@biloxi.com>;tag=a6c85cf\r\n")...)
	b = append(b, []byte("From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n")...)
	b = append(b, []byte("Call-ID: a84b4c76e66710\r\n")...)
	b = append(b, []byte("CSeq: 314159 INVITE\r\n")...)
	b = append(b, []byte("Contact: <sip:bob@192.0.2.4>\r\n")...)
	b = append(b, []byte("Content-Type: application/sdp\r\n")...)
	b = append(b, []byte("Content-Length: 131\r\n")...)

	return b
}

func (s *Server) buildRequestBYE() []byte {
	var b []byte
	b = append(b, []byte("BYE sip:alice@pc33.atlanta.com SIP/2.0")...)
	b = append(b, []byte("Via: SIP/2.0/UDP 192.0.2.4;branch=z9hG4bKnashds10")...)
	b = append(b, []byte("Max-Forwards: 70")...)
	b = append(b, []byte("From: Bob <sip:bob@biloxi.com>;tag=a6c85cf")...)
	b = append(b, []byte("To: Alice <sip:alice@atlanta.com>;tag=1928301774")...)
	b = append(b, []byte("Call-ID: a84b4c76e66710")...)
	b = append(b, []byte("CSeq: 231 BYE")...)
	b = append(b, []byte("Content-Length: 0")...)

	return b
}

var (
	space        = []byte(" ")
	end          = []byte("\r\n")
	methodINVITE = []byte("INVITE")
)

func isValidRequestINVITE(b []byte) bool {
	requestMsg := bytes.Split(b, end)
	// check request line
	if !isValidRequestLine(requestMsg[0], methodINVITE) {
		return false
	}

	// check request header
	if !isValidRequestHeader(requestMsg[1:]) {
		return false
	}

	return true
}

func isValidRequestLine(requestLine, method []byte) bool {
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
