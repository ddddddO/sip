package sip

type State string

// https://ja.wikipedia.org/wiki/Session_Initiation_Protocol [SIP における標準的なシーケンス]
const (
	StateINIT      = State("init") // session initialized
	StateINVITE    = State("invite")
	StateRINGING   = State("ringing")
	StateOK        = State("ok")
	StateACK       = State("ack")
	StateCONNECTED = State("connected") // session established
	StateBYE       = State("bye")
)
