package sip

import (
	"log"
)

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
	for raddr := range c.ssmap {
		switch c.ssmap[raddr].GetState() {
		case StateINIT:
			inviteReq := buildRequestINVITE()
			c.ssmap[raddr].Write(inviteReq)
			c.ssmap[raddr].ChangeState(StateINVITE)

			b, err := c.ssmap[raddr].Read()
			if err != nil {
				return err
			}
			log.Printf("debug\n%s", string(b))
		}
	}
	return nil
}
