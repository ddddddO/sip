package sip

import (
	"bytes"
)

func buildResponseRINGING() []byte {
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

func buildResponseOK() []byte {
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

// 一旦、180のみtrue
func isValidStatusCode1XX(b []byte) bool {
	responseMsg := bytes.Split(b, []byte("\r\n"))[0]
	splited := bytes.Split(responseMsg, []byte(" "))

	if !bytes.Equal(splited[1], []byte("180")) {
		return false
	}
	return true
}

// 一旦、200のみtrue
func isValidStatusCode2XX(b []byte) bool {
	responseMsg := bytes.Split(b, []byte("\r\n"))[0]
	splited := bytes.Split(responseMsg, []byte(" "))

	if !bytes.Equal(splited[1], []byte("200")) {
		return false
	}
	return true
}
