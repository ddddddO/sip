package sip

import (
	"bufio"
	"errors"
	"net"
)

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
