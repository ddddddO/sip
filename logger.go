package sip

import (
	"log"
)

func LogInfo(raddr string, ss *Session) {
	log.Printf("\nremote address: %s\nstatus: %s",
		raddr,
		ss.GetState())
}
