package server

import (
	"fmt"
	"net"
)

type MessageReceiver struct {
	conn *net.Conn
	ch   chan []byte
}

func NewMessageReceiver(conn *net.Conn) *MessageReceiver {
	return &MessageReceiver{
		conn: conn,
		ch:   make(chan []byte),
	}
}

func (mr *MessageReceiver) Start() {
	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := (*mr.conn).Read(buf)
			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("Connection closed")
					break
				}
				panic(err)
			}
			mr.ch <- buf[:n]
		}
	}()

	for {
		msg := mr.receive()
		fmt.Printf("Msg sent from %s: %s\n", (*mr.conn).RemoteAddr(), msg)
	}
}

func (mr *MessageReceiver) receive() []byte {
	return <-mr.ch
}
