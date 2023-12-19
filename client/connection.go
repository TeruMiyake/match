package client

import (
	"net"

	"github.com/TeruMiyake/match/message"
)

type ClientConnection struct {
	conn  *net.TCPConn
	msgCh chan message.Message
}

func NewClientConnection(conn *net.TCPConn) *ClientConnection {
	return &ClientConnection{
		conn:  conn,
		msgCh: make(chan message.Message),
	}
}

func (cc *ClientConnection) Start() chan<- message.Message {
	go func() {
		for {
			msg := <-cc.msgCh
			bytes := msg.Bytes()
			cc.conn.Write(bytes)
		}
	}()

	return cc.msgCh
}
