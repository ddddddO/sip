package sip

import (
	"bufio"
	"net"
)

//まず、クライアント・サーバ1対1で考える

type Client struct {
	conn *net.UDPConn // TODO: 接続先が複数の場合があるから、conn->connsになる？そうなるとかなり変更しそう
	br   *bufio.Reader
	bw   *bufio.Writer
}

func NewClient(raddr string) *Client {
	conn, err := net.Dial("udp4", raddr)
	if err != nil {
		panic(err)
	}
	udpConn, ok := conn.(*net.UDPConn)
	if !ok {
		panic("can't assertion")
	}

	return &Client{
		conn: udpConn,
		br:   bufio.NewReader(conn),
		bw:   bufio.NewWriter(conn),
	}
}

func (c *Client) Read() ([]byte, error) {
	size := c.br.Size()
	buf := make([]byte, size)

	_, err := c.br.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *Client) Write(b []byte) error {
	_, err := c.bw.Write(b)
	if err != nil {
		return err
	}
	if err := c.bw.Flush(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

type Server struct {
	conn *net.UDPConn
	br   *bufio.Reader
	bw   *bufio.Writer
}

func NewServer(laddr string) *Server {
	conn, err := net.ListenPacket("udp", laddr)
	if err != nil {
		panic(err)
	}
	udpConn, ok := conn.(*net.UDPConn)
	if !ok {
		panic("can't assertion")
	}

	return &Server{
		conn: udpConn,
		br:   bufio.NewReader(udpConn),
		bw:   bufio.NewWriter(udpConn),
	}
}

func (s *Server) Read() (net.Addr, []byte, error) {
	size := s.br.Size()
	b := make([]byte, size)
	_, raddr, err := s.conn.ReadFrom(b)
	if err != nil {
		return nil, nil, err
	}
	return raddr, b, nil
}

func (s *Server) Write(b []byte, raddr net.Addr) error {
	_, err := s.conn.WriteTo(b, raddr)
	if err != nil {
		return err
	}
	return nil
}
