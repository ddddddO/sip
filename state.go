package sip

type State string

const (
	StateINIT      = State("init")
	StateINVITE    = State("invite")
	StateOK        = State("ok")
	StateACK       = State("ack")
	StateCONNECTED = State("connected")
)
