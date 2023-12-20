package server

import (
	"fmt"

	"github.com/TeruMiyake/match/message"
)

type MessageReceiver struct {
	mre *message.MessageReader
	buf []byte
}

func NewMessageReceiver() *MessageReceiver {
	return &MessageReceiver{
		mre: message.NewMessageReader(),
		buf: make([]byte, 0, 1024),
	}
}

func (mr *MessageReceiver) Receive(b []byte) {
	mr.mre.Write(b)

	for msg, err := mr.mre.Read(); err == nil; msg, err = mr.mre.Read() {
		// TODO: Surrender の実装
		fmt.Printf("Msg received : %+v\n", msg)
	}
}
