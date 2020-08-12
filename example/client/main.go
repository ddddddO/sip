package main

import (
	"log"

	"github.com/ddddddO/sip"
)

func main() {
	log.Print("start sip uac(client proccess)")
	raddr := "localhost:5060"
	client := sip.NewClient(raddr)
	defer client.Close()

	log.Print("send to server")
	err := client.Write(buildRequestINVITE())
	if err != nil {
		panic(err)
	}

	buf, err := client.Read()
	if err != nil {
		panic(err)
	}
	log.Print(string(buf))
}

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
