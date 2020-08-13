package sip

import (
	"sync"
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
	errCh := make(chan error)

	wg := &sync.WaitGroup{}
	for raddr := range c.ssmap {
		wg.Add(1)
		go c.run(raddr, wg, errCh)

	}
	wg.Wait()

	if len(errCh) > 0 {
		return <-errCh
	}

	return nil
}

func (c *Client) run(raddr string, wg *sync.WaitGroup, errCh chan<- error) {
CONNECTED:
	for {
		switch c.ssmap[raddr].GetState() {
		case StateINIT:
			LogInfo(raddr, c.ssmap[raddr])
			inviteReq := buildRequestINVITE()
			c.ssmap[raddr].Write(inviteReq)
			c.ssmap[raddr].ChangeState(StateINVITE)

			res, err := c.ssmap[raddr].Read()
			if err != nil {
				errCh <- err
			}
			if !isValidStatusCode1XX(res) {
				// TODO: ...
			}
		case StateINVITE:
			res, err := c.ssmap[raddr].Read()
			if err != nil {
				errCh <- err
			}
			LogInfo(raddr, c.ssmap[raddr])
			if isValidStatusCode2XX(res) {
				ackReq := buildRequestACK()
				c.ssmap[raddr].Write(ackReq)
				c.ssmap[raddr].ChangeState(StateCONNECTED)
			} else {
				// TODO: ...
			}
		case StateCONNECTED:
			LogInfo(raddr, c.ssmap[raddr])
			break CONNECTED
		}
	}
	wg.Done()
}
