package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
)

type OpCode int

const (
	Ping OpCode = iota
	Close
	Text
	G
)

type WebSocket struct {
	Inbox  chan []byte
	Closed chan bool
	rw     *bufio.ReadWriter
}

func (ws *WebSocket) Worker() {

	for {

		select {
		case <-ws.Closed:
			return
		case message := <-ws.Inbox:
			ws.Write(message)
		}
	}

}

func (ws *WebSocket) Listen() {
	/*
		go func() {
			frame, code, err := ReadFrame(reader)
			if err != nil {
				fmt.Println(err.Error())
				//done <- true
				close(d)
				return
			} else if code == Close {

				close(d)
				return
			}
		}()
	*/
}

func (ws *WebSocket) Read() (string, OpCode, error) {
	return ReadFrame(ws.rw.Reader)
}

func ReadFrame(reader *bufio.Reader) (string, OpCode, error) {
	// | isfinal? | x x x | opcode(4) |
	// | ismask? | length(7) |
	// | mask (32) |
	header := make([]byte, 2)

	_, err := reader.Read(header)

	if err != nil {
		return "", G, err
	}

	//	fmt.Printf("header length read : %d \n", hlen)

	//var isFinal = header[0] >> 7
	var opcode = header[0] & 15
	//var isMasked = header[1] >> 7
	var length = int(header[1] & 127)
	//fmt.Printf("raw header : %b %b \n", header[0], header[1])
	//fmt.Printf("header : %d %d %d %d \n", isFinal, opcode, isMasked, length)

	if opcode == 8 {
		return "", Close, nil
	}

	//client to server always has a mask
	mask := make([]byte, 4)
	_, _ = reader.Read(mask)

	body := make([]byte, length)

	_, _ = reader.Read(body)

	for i := 0; i < length; i++ {
		/*
		   next,err := reader.ReadByte()
		   if err!=nil{
		     log.Printf("error reading frame: %v", err)
		     break
		   }*/
		//unmask
		body[i] = body[i] ^ mask[i%4]
	}

	s := string(body[:length])
	//fmt.Printf("string : %s \n", s)
	return s, Text, nil

}

func (ws *WebSocket) Write(b []byte) (n int, err error) {

	rw := ws.rw
	length := len(b)
	max16 := 65535
	//header := make([]byte, 0)
	//	fmt.Println(length)
	//as long as this works in all browsers, then its fine for < 126
	//len64 := (length >> 16) & 255
	blen := (length >> 16) & 255
	llen := (length >> 8) & 255
	rlen := length & 255

	header := []byte{129}

	if length > max16 {
		header = append(header, 127, 0, 0, 0, 0, 0, byte(blen))
	} else if length >= 126 {
		header = append(header, 126, byte(llen), byte(rlen))
	} else {
		header = append(header, byte(length))
	}
	rw.Write(header)

	//rw.Write([]byte{129, 127, byte(0), byte(0), byte(0), byte(len64), byte(llen), byte(rlen)})

	/*
	   if length > 125 {

	   } else {
	     rw.Write([]byte{129, byte(length)})
	   }
	*/
	rw.Write(b)
	rw.Flush()

	return length, nil

}

func upgrade(w http.ResponseWriter, req *http.Request) *WebSocket {

	//pid := req.URL.Query().Get("id")

	var guid = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	key := req.Header.Get("Sec-WebSocket-Key")

	hash := sha1.New()
	hash.Write([]byte(key))
	hash.Write(guid)

	shaed := hash.Sum(nil)
	var challengeresponse = base64.StdEncoding.EncodeToString(shaed)

	h, _ := w.(http.Hijacker)
	_, rw, _ := h.Hijack()
	//defer conn.Close()

	buf := make([]byte, 0, 4096)

	buf = append(buf, "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: "...)
	buf = append(buf, challengeresponse...)
	buf = append(buf, "\r\n"...)
	buf = append(buf, "\r\n"...)

	rw.Write(buf)
	rw.Flush()

	ws := &WebSocket{
		rw:     rw,
		Closed: make(chan bool),
		Inbox:  make(chan []byte, 10),
	}

	go ws.Worker()

	return ws

}
