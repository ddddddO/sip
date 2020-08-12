package sip

import (
	"bufio"
	"errors"
	"log"
	"net"
)

//まず、クライアント・サーバ1対1で考える

// Key: remote address / Value: Session
type SessionMap map[string]*Session

type Session struct {
	conn *net.UDPConn
	br   *bufio.Reader
	bw   *bufio.Writer

	state State
}

type State string

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

	return &Session{
		conn: udpConn,
		br:   bufio.NewReader(conn),
		bw:   bufio.NewWriter(conn),

		state: "init",
	}
}

func (ss *Session) Read() ([]byte, error) {
	if ss.br == nil {
		return nil, errors.New("can't access")
	}

	size := ss.br.Size()
	buf := make([]byte, size)

	_, err := ss.br.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (ss *Session) Write(b []byte) error {
	if ss.bw == nil {
		return errors.New("can't access")
	}
	_, err := ss.bw.Write(b)
	if err != nil {
		return err
	}
	if err := ss.bw.Flush(); err != nil {
		return err
	}
	return nil
}

func (ss *Session) Close() error {
	return ss.conn.Close()
}

type Client struct {
	ssmap SessionMap
}

func NewClient() *Client {
	return &Client{
		ssmap: SessionMap{},
	}
}

func (c *Client) AddSession(raddr string, session *Session) {
	if _, ok := c.ssmap[raddr]; !ok {
		c.ssmap[raddr] = session
	}
}

func (c *Client) Run() error {
	for key := range c.ssmap {
		c.ssmap[key].Write(c.buildRequestINVITE())

		b, err := c.ssmap[key].Read()
		if err != nil {
			return err
		}
		log.Print(string(b))
	}
	return nil
}

// https://tools.ietf.org/html/rfc3261#section-7.1
func (c *Client) buildRequestINVITE() []byte {
	var b []byte
	b = append(b, []byte("INVITE sip:bob@biloxi.com SIP/2.0\r\n")...)
	// NOTE: ヘッダーに\r\nは必要でよい？
	b = append(b, []byte("Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds\r\n")...)
	b = append(b, []byte("Max-Forwards: 70\r\n")...)
	b = append(b, []byte("To: Bob <sip:bob@biloxi.com>\r\n")...)
	b = append(b, []byte("From: Alice <sip:alice@atlanta.com>;tag=1928301774\r\n")...)
	b = append(b, []byte("Call-ID: a84b4c76e66710@pc33.atlanta.com\r\n")...)
	b = append(b, []byte("CSeq: 314159 INVITE\r\n")...)
	b = append(b, []byte("Contact: <sip:alice@pc33.atlanta.com>\r\n")...)
	b = append(b, []byte("Content-Type: application/sdp\r\n")...)
	b = append(b, []byte("Content-Length: 142\r\n")...)

	return b
}

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
			state: "init",
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
