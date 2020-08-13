package sip

import (
	"bytes"
)

// https://tools.ietf.org/html/rfc3261#section-7.1
func buildRequestINVITE() []byte {
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

func buildRequestACK() []byte {
	var b []byte
	b = append(b, []byte("ACK sip:bob@192.0.2.4 SIP/2.0")...)
	b = append(b, []byte("Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bKnashds9")...)
	b = append(b, []byte("Max-Forwards: 70")...)
	b = append(b, []byte("To: Bob <sip:bob@biloxi.com>;tag=a6c85cf")...)
	b = append(b, []byte("From: Alice <sip:alice@atlanta.com>;tag=1928301774")...)
	b = append(b, []byte("Call-ID: a84b4c76e66710")...)
	b = append(b, []byte("CSeq: 314159 ACK")...)
	b = append(b, []byte("Content-Length: 0")...)

	return b
}

func buildRequestBYE() []byte {
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
